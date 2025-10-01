package application

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/newinternal/application/dto"
	applicationports "github.com/jairoprogramador/fastdeploy/newinternal/application/ports"
	deploymentaggregates "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	deploymentports "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/ports"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	executionstateaggregates "github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/aggregates"
	executionstateports "github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/ports"
	executionstateservices "github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/services"
	executionstatevos "github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/vos"
	orchestrationaggregates "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/aggregates"
	orchestrationports "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/ports"
	orchestrationservices "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/services"
	orchestrationvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git"
)

// OrchestrationService es el servicio de aplicación que coordina el caso de uso
// principal: la ejecución de una orden de despliegue.
type OrchestrationService struct {
	templateRepo deploymentports.TemplateRepository
	orderRepo    orchestrationports.OrderRepository
	historyRepo  executionstateports.HistoryRepository
	varResolver  orchestrationservices.VariableResolver
	fpService    executionstateservices.FingerprintService
	workspaceMgr applicationports.WorkspaceManager
	cmdExecutor  applicationports.CommandExecutor
}

func NewOrchestrationService(
	templateRepo deploymentports.TemplateRepository,
	orderRepo orchestrationports.OrderRepository,
	historyRepo executionstateports.HistoryRepository,
	varResolver orchestrationservices.VariableResolver,
	fpService executionstateservices.FingerprintService,
	workspaceMgr applicationports.WorkspaceManager,
	cmdExecutor applicationports.CommandExecutor,
) *OrchestrationService {
	return &OrchestrationService{
		templateRepo: templateRepo,
		orderRepo:    orderRepo,
		historyRepo:  historyRepo,
		varResolver:  varResolver,
		fpService:    fpService,
		workspaceMgr: workspaceMgr,
		cmdExecutor:  cmdExecutor,
	}
}

func (s *OrchestrationService) ExecuteOrder(req dto.ExecuteOrderRequest) (*orchestrationaggregates.Order, error) {
	// 1. Obtener la "receta" Y la ruta local del repo.
	template, templateRepoPath, err := s.templateRepo.GetTemplate(req.Ctx, req.TemplateSource)
	if err != nil {
		return nil, err
	}
	fmt.Println("templateRepoPath", templateRepoPath)

	repositoryName, err := s.templateRepo.GetRepositoryName(req.TemplateSource.RepoURL())
	if err != nil {
		return nil, err
	}
	fmt.Println("repositoryName", repositoryName)

	// 2. Determinar el ambiente objetivo.
	targetEnv, err := s.getEnvironment(template.Environments(), req.EnvironmentName)
	if err != nil {
		return nil, err
	}
	fmt.Println("targetEnv", targetEnv)

	stepsToExecute := s.getStepsToExecute(template, req.FinalStepName)

	stepVarRepo := git.NewVariableRepository(templateRepoPath)

	allVars := req.InitialVariables
	for _, stepDef := range stepsToExecute {
		stepVars, err := stepVarRepo.Load(req.Ctx, targetEnv, stepDef)
		if err != nil {
			return nil, fmt.Errorf("error al cargar las variables para el paso '%s': %w", stepDef.Name(), err)
		}
		allVars = append(allVars, stepVars...)
	}

	// 4. Crear el agregado Order. El dominio se encarga de la lógica de negocio.
	order, err := orchestrationaggregates.NewOrder(
		orchestrationvos.NewOrderID(),
		template,
		targetEnv,
		req.FinalStepName,
		req.SkippedStepNames,
		allVars,
	)
	if err != nil {
		return nil, fmt.Errorf("error al crear la orden de ejecución: %w", err)
	}

	// 5. Bucle de ejecución de la orden (pasos y comandos).
	fmt.Printf("Iniciando despliegue hasta el paso '%s' en el ambiente '%s'...\n", req.FinalStepName, targetEnv.Name())

	// Cargar los últimos recibos exitosos de test y supply
	latestTestReceipt, _ := s.historyRepo.Find(req.Ctx, "test")
	latestSupplyReceipt, _ := s.historyRepo.Find(req.Ctx, "supply")

	for _, stepExec := range order.StepExecutions() {
		if stepExec.Status() == orchestrationvos.StepStatusSkipped {
			fmt.Printf("--- Omitiendo paso: %s ---\n", stepExec.Name())
			continue
		}

		stepDef := findStepDefinition(template.Steps(), stepExec.Name())

		// --- LÓGICA DE DECISIÓN DE CACHÉ ---
		workspacePath, err := s.workspaceMgr.PrepareStepWorkspace(req.Ctx, req.ProjectName, targetEnv, stepDef, repositoryName)
		if err != nil {
			return order, fmt.Errorf("error al preparar el workspace para el paso '%s': %w", stepDef.Name(), err)
		}

		currentCodeFp, err := s.fpService.CalculateCodeFingerprint(req.Ctx, req.ProjectRootPath, []string{".fdignore"})
		if err != nil {
			return order, fmt.Errorf("error al calcular el fingerprint del código: %w", err)
		}

		currentEnvFp, err := s.fpService.CalculateEnvironmentFingerprint(req.Ctx, workspacePath)
		if err != nil {
			return order, fmt.Errorf("error al calcular el fingerprint del ambiente: %w", err)
		}

		if s.shouldSkipStep(stepDef, currentCodeFp, currentEnvFp, latestTestReceipt, latestSupplyReceipt) {
			fmt.Printf("--- Omitiendo paso (ya ejecutado con éxito): %s ---\n", stepExec.Name())
			order.MarkStepAsCached(stepExec.Name()) // <-- NUEVO MÉTODO EN EL AGREGADO
			continue
		}

		fmt.Printf("--- Ejecutando paso: %s ---\n", stepExec.Name())
		for _, cmdExec := range stepExec.CommandExecutions() {
			fmt.Printf("-> Ejecutando comando: %s\n", cmdExec.Name())

			// Interpolar el comando antes de ejecutarlo.
			interpolatedCmd, err := s.varResolver.Interpolate(cmdExec.Definition().CmdTemplate(), order.VariableMap())
			if err != nil {
				// Este error es crítico, una variable no fue encontrada.
				return order, fmt.Errorf("error al interpolar el comando '%s': %w", cmdExec.Name(), err)
			}

			// Ejecutar el comando.
			log, exitCode, err := s.cmdExecutor.Execute(req.Ctx, cmdExec.Definition().Workdir(), interpolatedCmd)
			if err != nil {
				return order, fmt.Errorf("error del sistema al ejecutar el comando '%s': %w", cmdExec.Name(), err)
			}

			// Actualizar el estado del agregado con el resultado.
			err = order.MarkCommandAsCompleted(stepExec.Name(), cmdExec.Name(), interpolatedCmd, log, exitCode, s.varResolver)
			if err != nil {
				return order, err
			}

			// Persistir el estado después de cada comando.
			if err := s.orderRepo.Save(req.Ctx, order, req.ProjectName); err != nil {
				return order, err
			}

			// Si el estado de la orden es Failed, detener toda la ejecución.
			if order.Status() == orchestrationvos.OrderStatusFailed {
				fmt.Printf("❌ La orden ha fallado en el paso '%s', comando '%s'. Abortando. !!!\n", stepExec.Name(), cmdExec.Name())
				return order, nil // La orden falló, pero la operación del servicio fue "exitosa" en su manejo.
			}
		}

		// --- LÓGICA DE GUARDADO DE RECIBO ---
		if stepExec.Status() == orchestrationvos.StepStatusSuccessful {
			if stepExec.Name() == "test" {
				receipt, _ := executionstateaggregates.NewExecutionReceipt("test", currentCodeFp, executionstatevos.Fingerprint{}, order.ID())
				latestTestReceipt.AddReceipt(receipt)
				s.historyRepo.Save(req.Ctx, latestTestReceipt)
			}
			if stepExec.Name() == "supply" {
				receipt, _ := executionstateaggregates.NewExecutionReceipt("supply", executionstatevos.Fingerprint{}, currentEnvFp, order.ID())
				latestSupplyReceipt.AddReceipt(receipt)
				s.historyRepo.Save(req.Ctx, latestSupplyReceipt)
			}
		}
	}

	fmt.Println("🎉 Ejecución de la orden completada con éxito.")
	return order, nil
}

func (s *OrchestrationService) shouldSkipStep(
	stepDef deploymententities.StepDefinition,
	currentCodeFp, currentEnvFp executionstatevos.Fingerprint,
	latestTestHistory *executionstateaggregates.StepExecutionHistory,
	latestSupplyHistory *executionstateaggregates.StepExecutionHistory,
) bool {
	verifications := stepDef.VerificationTypes()
	if len(verifications) == 0 {
		return false // Sin verificaciones, siempre se ejecuta.
	}

	for _, verification := range verifications {
		if verification == deploymentvos.VerificationTypeCode {
			if latestTestHistory == nil || latestTestHistory.FindMatch(currentCodeFp, executionstatevos.Fingerprint{}) == nil {
				return false // El código ha cambiado o nunca se ha ejecutado.
			}
		}
		if verification == deploymentvos.VerificationTypeEnv {
			if latestSupplyHistory == nil || latestSupplyHistory.FindMatch(executionstatevos.Fingerprint{}, currentEnvFp) == nil {
				return false // El ambiente ha cambiado o nunca se ha ejecutado.
			}
		}
	}

	return true // Todas las verificaciones pasaron, se puede omitir.
}

// findEnvironment es una función de ayuda para buscar un ambiente por su nombre o valor.
func findEnvironment(envs []deploymentvos.Environment, nameOrValue string) (deploymentvos.Environment, error) {
	for _, env := range envs {
		if env.Name() == nameOrValue || env.Value() == nameOrValue {
			return env, nil
		}
	}
	return deploymentvos.Environment{}, fmt.Errorf("el ambiente '%s' no se encontró en la definición de la plantilla", nameOrValue)
}

func (s *OrchestrationService) getEnvironment(
	environments []deploymentvos.Environment, envName string) (deploymentvos.Environment, error) {

	var targetEnv deploymentvos.Environment
	if envName == "" {
		if len(environments) == 0 {
			return deploymentvos.Environment{}, fmt.Errorf("no hay ambientes configurados en la plantilla; se debe especificar al menos uno")
		}
		targetEnv = environments[0]
	} else {
		env, err := findEnvironment(environments, envName)
		if err != nil {
			return deploymentvos.Environment{}, err
		}
		targetEnv = env
	}
	return targetEnv, nil
}

// getStepsToExecute devuelve la secuencia de pasos a ejecutar según el paso final.
func (s *OrchestrationService) getStepsToExecute(template *deploymentaggregates.DeploymentTemplate, finalStepName string) []deploymententities.StepDefinition {
	var stepsToExecute []deploymententities.StepDefinition
	for _, step := range template.Steps() {
		stepsToExecute = append(stepsToExecute, step)
		if step.Name() == finalStepName {
			break
		}
	}
	return stepsToExecute
}

func findStepDefinition(steps []deploymententities.StepDefinition, stepName string) deploymententities.StepDefinition {
	for _, step := range steps {
		if step.Name() == stepName {
			return step
		}
	}
	return deploymententities.StepDefinition{} // Should not happen if stepName is valid
}
