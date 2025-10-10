package application

import (
	"fmt"
	"errors"

	"github.com/jairoprogramador/fastdeploy/internal/application/dto"
	applicationports "github.com/jairoprogramador/fastdeploy/internal/application/ports"
	deploymententities "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entities"
	deploymentvos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	domaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	executionstateaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
	executionstateports "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/ports"
	executionstateservices "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/services"
	executionstatevos "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"
	orchestrationaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/aggregates"
	orchestrationentities "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/entities"
	orchestrationports "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/ports"
	orchestrationservices "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/services"
	orchestrationvos "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

type OrchestrationService struct {
	stepVariableRepo orchestrationports.StepVariableRepository
	orderRepo        orchestrationports.OrderRepository
	scopeRepo        executionstateports.ScopeRepository
	varResolver      orchestrationservices.VariableResolver
	fpService        executionstateservices.FingerprintService
	workspaceMgr     applicationports.WorkspaceManager
	cmdExecutor      applicationports.CommandExecutor
	varsRepo         executionstateports.VarsRepository
	stateRepo        executionstateports.StateRepository
}

func NewOrchestrationService(
	stepVariableRepo orchestrationports.StepVariableRepository,
	orderRepo orchestrationports.OrderRepository,
	scopeRepo executionstateports.ScopeRepository,
	varResolver orchestrationservices.VariableResolver,
	fpService executionstateservices.FingerprintService,
	workspaceMgr applicationports.WorkspaceManager,
	cmdExecutor applicationports.CommandExecutor,
	varsRepo executionstateports.VarsRepository,
	stateRepo executionstateports.StateRepository,
) *OrchestrationService {
	return &OrchestrationService{
		stepVariableRepo: stepVariableRepo,
		orderRepo:        orderRepo,
		scopeRepo:        scopeRepo,
		varResolver:      varResolver,
		fpService:        fpService,
		workspaceMgr:     workspaceMgr,
		cmdExecutor:      cmdExecutor,
		varsRepo:         varsRepo,
		stateRepo:        stateRepo,
	}
}

func (s *OrchestrationService) ExecuteOrder(req dto.OrderRequest) (*orchestrationaggregates.Order, error) {

	allVars, err := s.getVariablesInit(req)
	if err != nil {
		return nil, err
	}

	statusLatestSteps, err := s.stateRepo.FindStepStatus()
	if err != nil {
		return nil, err
	}

	stateCurrentCode, err := s.fpService.CalculateCodeFingerprint()
	if err != nil {
		return nil, err
	}

	latestCodeStatusHistory, err := s.scopeRepo.FindCodeStateHistory()
	if err != nil {
		return nil, err
	}

	latestEnvironmentStatusHistoryMap, err := s.findEnvironmentStateHistory(req)
	if err != nil {
		return nil, err
	}

	order, err := orchestrationaggregates.NewOrder(
		orchestrationvos.NewOrderID(),
		req.Template,
		req.Environment,
		req.FinalStep,
		req.SkippedStepNames,
		allVars,
	)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Iniciando despliegue hasta el paso '%s' en el ambiente '%s'...\n", req.FinalStep, req.Environment.Name())
	for _, stepExec := range order.StepExecutions() {
		if stepExec.Status() == orchestrationvos.StepStatusSkipped {
			fmt.Println("\n-----------------------------------------------")
			fmt.Printf("------------- OMITIENDO STEP: %s -------------\n", stepExec.Name())
			fmt.Println("-----------------------------------------------")
			continue
		}

		stateCurrentEnvironmentStep, err := s.fpService.CalculateStepFingerprint(stepExec.Name())
		if err != nil {
			return order, err
		}

		if statusLatestSteps.IsStepAlreadyExecuted(stepExec.Name()) {
			stepDef := s.findStepDefinition(req.Template.Steps(), stepExec.Name())
			if s.thereAreChanges(
				stepDef.VerificationTypes(),
				stepExec.Name(),
				stateCurrentCode,
				stateCurrentEnvironmentStep,
				latestCodeStatusHistory,
				latestEnvironmentStatusHistoryMap) {

				err = s.processStep(req, order, stepExec, stateCurrentCode, stateCurrentEnvironmentStep)
				if err != nil {
					return order, err
				}
			} else {
				order.MarkStepAsCached(stepExec.Name())
				fmt.Printf("---------------- OMITIENDO STEP: %s -----------------\n", stepExec.Name())
				fmt.Println("-------------------------------------------------------")
				continue
			}
		} else {
			err = s.processStep(req, order, stepExec, stateCurrentCode, stateCurrentEnvironmentStep)
			if err != nil {
				return order, err
			}
		}
	}

	if order.Status() == orchestrationvos.OrderStatusFailed {
		fmt.Println("❌ La ejecución de la orden ha fallado")
	}
	if order.Status() == orchestrationvos.OrderStatusSuccessful {
		fmt.Println("🎉 Ejecución de la orden completada con éxito.")
	}
	return order, nil
}

func (s *OrchestrationService) processStep(
	req dto.OrderRequest,
	order *orchestrationaggregates.Order,
	stepExec *orchestrationentities.StepExecution,
	stateCurrentCode executionstatevos.Fingerprint,
	stateCurrentEnvironmentStep executionstatevos.Fingerprint) error {

	_, err := s.loadStepVariables(stepExec.Name(), order)
	if err != nil {
		return err
	}

	err = s.executeStep(req, order, stepExec)
	if err != nil {
		return err
	}

	err = s.saveVarsStore(order.VariableMap())
	if err != nil {
		return err
	}
	err = s.saveStateSteps(order.StepExecutions())
	if err != nil {
		return err
	}
	err = s.saveFingerprints(stepExec.Name(), stateCurrentCode, stateCurrentEnvironmentStep)
	if err != nil {
		return err
	}

	if err := s.orderRepo.Save(order, req.ProjectDom.Project().Name()); err != nil {
		return err
	}

	if stepExec.Status() == orchestrationvos.StepStatusFailed {
		msm := fmt.Sprintf("❌ La ejecución del paso '%s' ha fallado", stepExec.Name())
		return errors.New(msm)
	}
	if stepExec.Status() == orchestrationvos.StepStatusSuccessful {
		fmt.Printf("🎉 Ejecución del paso '%s' completada con éxito.\n", stepExec.Name())
	}
	return nil
}

func (s *OrchestrationService) executeStep(
	req dto.OrderRequest,
	order *orchestrationaggregates.Order,
	stepExec *orchestrationentities.StepExecution) error {

	workdirStep, err := s.workspaceMgr.Prepare(stepExec.Name())

	if err != nil {
		return err
	}
	order.AddVariable("step_workdir", workdirStep)

	fmt.Printf("--------------- EJECUTANDO STEP: %s ---------------\n", stepExec.Name())
	fmt.Println("--------------------------------------------------------")
	for _, cmdExec := range stepExec.CommandExecutions() {
		fmt.Printf("-> Ejecutando comando: %s\n", cmdExec.Name())

		workdirCmd := cmdExec.Definition().Workdir()
		if workdirCmd != "" && workdirStep != "" {
			workdirCmd = s.cmdExecutor.CreateWorkDir(workdirStep, workdirCmd)
			order.AddVariable("command_workdir", workdirCmd)
		}

		interpolatedCmd, err := s.varResolver.Interpolate(cmdExec.Definition().CmdTemplate(), order.VariableMap())
		if err != nil {
			return err
		}

		for _, templatePath := range cmdExec.Definition().TemplateFiles() {
			pathTemplateFile := s.cmdExecutor.CreateWorkDir(workdirCmd, templatePath)

			err = s.varResolver.ProcessTemplate(pathTemplateFile, order.VariableMap())
			if err != nil {
				return err
			}
		}

		log, exitCode, err := s.cmdExecutor.Execute(req.Ctx, workdirCmd, interpolatedCmd)
		if err != nil {
			return err
		}

		err = order.MarkCommandAsCompleted(stepExec.Name(), cmdExec.Name(), interpolatedCmd, log, exitCode, s.varResolver)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OrchestrationService) saveFingerprints(
	stepName string,
	stateCurrentCode executionstatevos.Fingerprint,
	currentEnvFp executionstatevos.Fingerprint) error {

	if stepName == "test" {
		receiptCode, err := executionstateaggregates.NewScopeCodeReceipt(stateCurrentCode)
		if err != nil {
			return err
		}

		stateLatestCodeHistory, err := s.scopeRepo.FindCodeStateHistory()
		if err != nil {
			return err
		}

		stateLatestCodeHistory.AddReceipt(receiptCode)
		return s.scopeRepo.SaveCodeStateHistory(stateLatestCodeHistory)
	} else {
		receiptEnv, err := executionstateaggregates.NewScopeEnvironmentReceipt(currentEnvFp)
		if err != nil {
			return err
		}
		scopeLatestEnvReceipts, err := s.scopeRepo.FindStepStateHistory(stepName)
		if err != nil {
			return err
		}

		scopeLatestEnvReceipts.AddReceipt(receiptEnv)
		return s.scopeRepo.SaveStepStateHistory(scopeLatestEnvReceipts, stepName)
	}
}

func (s *OrchestrationService) saveStateSteps(steps []*orchestrationentities.StepExecution) error {
	stateSteps := executionstateaggregates.NewStateSteps()
	for _, stepExec := range steps {
		successful := (stepExec.Status() == orchestrationvos.StepStatusSuccessful) || (stepExec.Status() == orchestrationvos.StepStatusCached)
		step, err := executionstatevos.NewStateStep(stepExec.Name(), successful)
		if err != nil {
			return err
		}
		stateSteps.AddStep(step)
	}
	return s.stateRepo.SaveStepStatus(stateSteps)
}

func (s *OrchestrationService) saveVarsStore(varsMap map[string]orchestrationvos.Variable) error {
	varsStore := []orchestrationvos.Variable{}
	for _, variable := range varsMap {
		varsStore = append(varsStore, variable)
	}
	err := s.varsRepo.Save(varsStore)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrchestrationService) findStepDefinition(steps []deploymententities.StepDefinition, stepName string) deploymententities.StepDefinition {
	for _, step := range steps {
		if step.Name() == stepName {
			return step
		}
	}
	return deploymententities.StepDefinition{}
}

func (s *OrchestrationService) getVariablesInit(req dto.OrderRequest) ([]orchestrationvos.Variable, error) {
	environment := req.Environment.Value()

	allVarsInit := []orchestrationvos.Variable{}

	storeVars, err := s.varsRepo.FindAll()
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, storeVars...)

	varsProject, err := s.getVariablesProject(req.ProjectDom)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, varsProject...)

	env, err := orchestrationvos.NewVariable("environment", environment)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, env)

	projectPath, err := orchestrationvos.NewVariable("project_workdir", req.ProjectPath)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, projectPath)

	toolName, err := orchestrationvos.NewVariable("tool_name", "fastdeploy")
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, toolName)

	return allVarsInit, nil
}

func (s *OrchestrationService) getVariablesProject(domModel *domaggregates.DeploymentObjectModel) ([]orchestrationvos.Variable, error) {
	varsProject := []orchestrationvos.Variable{}

	projectId, err := orchestrationvos.NewVariable("project_id", domModel.Project().IdString()[:8])
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectId)

	projectName, err := orchestrationvos.NewVariable("project_name", domModel.Project().Name())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectName)

	projectTeam, err := orchestrationvos.NewVariable("project_team", domModel.Project().Team())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectTeam)

	projectVersion, err := orchestrationvos.NewVariable("project_version", domModel.Project().Version())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectVersion)

	projectRevision, err := orchestrationvos.NewVariable("project_revision", domModel.Project().Revision())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectRevision)

	productId, err := orchestrationvos.NewVariable("product_id", domModel.Product().IdString()[:8])
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, productId)

	productName, err := orchestrationvos.NewVariable("product_name", domModel.Product().Name())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, productName)

	productTeam, err := orchestrationvos.NewVariable("product_team", domModel.Product().Team())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, productTeam)

	productOrganization, err := orchestrationvos.NewVariable("product_organization", domModel.Product().Organization())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, productOrganization)

	technologyType, err := orchestrationvos.NewVariable("technology_type", domModel.Technology().TypeTechnology())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, technologyType)

	technologySolution, err := orchestrationvos.NewVariable("technology_solution", domModel.Technology().Solution())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, technologySolution)

	technologyStack, err := orchestrationvos.NewVariable("technology_stack", domModel.Technology().Stack())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, technologyStack)

	technologyInfrastructure, err := orchestrationvos.NewVariable("technology_infrastructure", domModel.Technology().Infrastructure())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, technologyInfrastructure)

	return varsProject, nil
}

func (s *OrchestrationService) findEnvironmentStateHistory(req dto.OrderRequest) (map[string]*executionstateaggregates.ScopeReceiptHistory, error) {
	latestEnvironmentStatusHistoryMap := make(map[string]*executionstateaggregates.ScopeReceiptHistory)

	for _, stepExec := range req.Template.Steps() {
		if _, ok := req.SkippedStepNames[stepExec.Name()]; ok {
			continue
		}
		latestEnvironmentStatusHistory, err := s.scopeRepo.FindStepStateHistory(stepExec.Name())
		if err != nil {
			return nil, err
		}
		latestEnvironmentStatusHistoryMap[stepExec.Name()] = latestEnvironmentStatusHistory
	}

	return latestEnvironmentStatusHistoryMap, nil
}

func (s *OrchestrationService) thereAreChanges(
	verifications []deploymentvos.VerificationType,
	stepName string,
	stateCurrentCode, stateCurrentEnvironmentStep executionstatevos.Fingerprint,
	stateLatestCodeHistory *executionstateaggregates.ScopeReceiptHistory,
	stateLatestEnvironmentHistoryMap map[string]*executionstateaggregates.ScopeReceiptHistory) bool {

	if len(verifications) == 0 {
		return true
	}

	for _, verification := range verifications {
		if verification == deploymentvos.VerificationTypeCode {
			if stateLatestCodeHistory == nil || stateLatestCodeHistory.FindMatchCode(stateCurrentCode) == nil {
				fmt.Println("\n-------------------------------------------------------")
				fmt.Printf("------------- HAY CAMBIOS EN EL CODIGO -------------\n")
				return true
			} else {
				fmt.Println("\n-------------------------------------------------------")
				fmt.Printf("------------- NO HAY CAMBIOS EN EL CODIGO -------------\n")
			}
		}
		if verification == deploymentvos.VerificationTypeEnv {
			stateLatestEnvironmentHistory, ok := stateLatestEnvironmentHistoryMap[stepName]
			if !ok || stateLatestEnvironmentHistory == nil || stateLatestEnvironmentHistory.FindMatchEnvironment(stateCurrentEnvironmentStep) == nil {
				fmt.Println("\n-------------------------------------------------------")
				fmt.Printf("------------- HAY CAMBIOS EN EL AMBIENTE -------------\n")
				return true
			} else {
				fmt.Println("\n-------------------------------------------------------")
				fmt.Printf("------------- NO HAY CAMBIOS EN EL AMBIENTE -------------\n")
			}
		}
	}
	return false
}

func (s *OrchestrationService) loadStepVariables(stepName string, order *orchestrationaggregates.Order) ([]orchestrationvos.Variable, error) {
	loadedStepVars, err := s.stepVariableRepo.Load(stepName)
	if err != nil {
		return nil, err
	}

	stepVars := []orchestrationvos.Variable{}

	for _, stepVar := range loadedStepVars {
		interpolatedVar, err := s.varResolver.Interpolate(stepVar.Value(), order.VariableMap())
		if err != nil {
			return nil, err
		}

		variable, err := orchestrationvos.NewVariable(stepVar.Key(), interpolatedVar)
		if err != nil {
			return nil, err
		}

		order.AddVariable(stepVar.Key(), interpolatedVar)
		stepVars = append(stepVars, variable)
	}

	return stepVars, nil
}
