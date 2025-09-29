package application

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/newinternal/application/dto"
	"github.com/jairoprogramador/fastdeploy/newinternal/application/ports"
	deployment "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	deploymentports "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/ports"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/aggregates"
	orchestrationports "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/ports"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/services"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// OrchestrationService es el servicio de aplicación que coordina el caso de uso
// principal: la ejecución de una orden de despliegue.
type OrchestrationService struct {
	templateRepo deploymentports.TemplateRepository
	orderRepo    orchestrationports.OrderRepository
	stepVarRepo  orchestrationports.StepVariableRepository
	varResolver  services.VariableResolver
	cmdExecutor  ports.CommandExecutor
}

// NewOrchestrationService crea una nueva instancia de OrchestrationService.
// Las dependencias se inyectan siguiendo el principio de Inversión de Dependencias.
func NewOrchestrationService(
	templateRepo deploymentports.TemplateRepository,
	orderRepo orchestrationports.OrderRepository,
	stepVarRepo orchestrationports.StepVariableRepository,
	varResolver services.VariableResolver,
	cmdExecutor ports.CommandExecutor,
) *OrchestrationService {
	return &OrchestrationService{
		templateRepo: templateRepo,
		orderRepo:    orderRepo,
		stepVarRepo:  stepVarRepo,
		varResolver:  varResolver,
		cmdExecutor:  cmdExecutor,
	}
}

// ExecuteOrder orquesta todo el proceso de ejecución de un despliegue.
func (s *OrchestrationService) ExecuteOrder(req dto.ExecuteOrderRequest) (*aggregates.Order, error) {
	// 1. Obtener la "receta" del despliegue.
	template, err := s.templateRepo.GetTemplate(req.Ctx, req.TemplateSource)
	if err != nil {
		return nil, err
	}

	repositoryName, err := s.templateRepo.GetRepositoryName(req.TemplateSource.RepoURL())
	if err != nil {
		return nil, err
	}

	// 2. Determinar el ambiente objetivo.
	var targetEnv deploymentvos.Environment
	allEnvs := template.Environments()
	if req.EnvironmentName == "" {
		if len(allEnvs) == 0 {
			return nil, fmt.Errorf("no hay ambientes configurados en la plantilla; se debe especificar al menos uno")
		}
		targetEnv = allEnvs[0]
	} else {
		env, err := findEnvironment(allEnvs, req.EnvironmentName)
		if err != nil {
			return nil, err
		}
		targetEnv = env
	}

	// 3. Determinar los pasos a ejecutar y cargar todas sus variables.
	stepsToExecute := getStepsToExecute(template, req.FinalStepName)
	allVars := req.InitialVariables
	for _, stepDef := range stepsToExecute {
		stepVars, err := s.stepVarRepo.Load(req.Ctx, repositoryName, targetEnv, stepDef)
		if err != nil {
			return nil, fmt.Errorf("error al cargar las variables para el paso '%s': %w", stepDef.Name(), err)
		}
		allVars = append(allVars, stepVars...)
	}

	// 4. Crear el agregado Order. El dominio se encarga de la lógica de negocio.
	order, err := aggregates.NewOrder(
		vos.NewOrderID(),
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
	for _, stepExec := range order.StepExecutions() {
		if stepExec.Status() == vos.StepStatusSkipped {
			fmt.Printf("--- Omitiendo paso: %s ---\n", stepExec.Name())
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
			if order.Status() == vos.OrderStatusFailed {
				fmt.Printf("❌ La orden ha fallado en el paso '%s', comando '%s'. Abortando. !!!\n", stepExec.Name(), cmdExec.Name())
				return order, nil // La orden falló, pero la operación del servicio fue "exitosa" en su manejo.
			}
		}
	}

	fmt.Println("🎉 Ejecución de la orden completada con éxito.")
	return order, nil
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

// getStepsToExecute devuelve la secuencia de pasos a ejecutar según el paso final.
func getStepsToExecute(template *deployment.DeploymentTemplate, finalStepName string) []deploymententities.StepDefinition {
	var stepsToExecute []deploymententities.StepDefinition
	allSteps := template.Steps()
	for _, step := range allSteps {
		stepsToExecute = append(stepsToExecute, step)
		if step.Name() == finalStepName {
			break
		}
	}
	return stepsToExecute
}
