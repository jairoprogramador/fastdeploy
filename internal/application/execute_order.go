package application

import (
	"errors"
	"fmt"

	appDto "github.com/jairoprogramador/fastdeploy/internal/application/dto"
	appPor "github.com/jairoprogramador/fastdeploy/internal/application/ports"

	stateAgg "github.com/jairoprogramador/fastdeploy/internal/domain/state/aggregates"
	statePor "github.com/jairoprogramador/fastdeploy/internal/domain/state/ports"
	stateSer "github.com/jairoprogramador/fastdeploy/internal/domain/state/services"
	statevos "github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"

	orchAgg "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/aggregates"
	orchEnt "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/entities"
	orchPor "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/ports"
	orchSer "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/services"
	orchVos "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"

	depEnt "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entities"
	domAgg "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
)

type ExecuteOrder struct {
	orderRepo           orchPor.OrderRepository
	varResolver         orchSer.TemplateResolver
	fingerprintService  stateSer.FingerprintService
	workspaceMgr        appPor.WorkspaceManager
	cmdExecutor         appPor.CommandExecutor
	variablesRepository statePor.VariablesRepository
	stateRepository     statePor.ExecutionStateRepository
	statePolicyService  stateSer.ExecutionPolicyService
}

func NewExecuteOrder(
	orderRepo orchPor.OrderRepository,
	varResolver orchSer.TemplateResolver,
	fingerprintService stateSer.FingerprintService,
	workspaceMgr appPor.WorkspaceManager,
	cmdExecutor appPor.CommandExecutor,
	variablesRepository statePor.VariablesRepository,
	stateRepository statePor.ExecutionStateRepository,
	statePolicyService stateSer.ExecutionPolicyService,
) *ExecuteOrder {
	return &ExecuteOrder{
		orderRepo:           orderRepo,
		varResolver:         varResolver,
		fingerprintService:  fingerprintService,
		workspaceMgr:        workspaceMgr,
		cmdExecutor:         cmdExecutor,
		variablesRepository: variablesRepository,
		stateRepository:     stateRepository,
		statePolicyService:  statePolicyService,
	}
}

func (s *ExecuteOrder) Run(req appDto.OrderRequest) (*orchAgg.Order, error) {

	allVars, err := s.getVariablesInit(req)
	if err != nil {
		return nil, err
	}

	fingerprintCurrentCode, err := s.fingerprintService.GenerateFromSource(req.ProjectPath)
	if err != nil {
		return nil, err
	}

	order, err := orchAgg.NewOrder(
		orchVos.NewOrderID(),
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
	for _, stepExec := range order.StepsRecord() {
		if stepExec.Status() == orchVos.StepStatusSkipped {
			fmt.Println("\n-----------------------------------------------")
			fmt.Printf("------------- OMITIENDO STEP: %s -------------\n", stepExec.Name())
			fmt.Println("-----------------------------------------------")
			continue
		}

		stepDef := s.findStepDefinition(req.Template.Steps(), stepExec.Name())

		err = s.processStepVariables(stepExec.Name(), stepDef, order)
		if err != nil {
			return order, err
		}
		varsMap := make(map[string]string)
		for key, variable := range order.Outputs() {
			varsMap[key] = variable.Value()
		}

		fingerprintCurrentVars, err := s.fingerprintService.GenerateFromStepVariables(varsMap)
		if err != nil {
			return order, err
		}

		fingerprintCurrentStep, err := s.fingerprintService.GenerateFromStepDefinition(req.TemplatePath, stepExec.Name())
		if err != nil {
			return order, err
		}

		stateLatest, err := s.stateRepository.FindByStepName(stepExec.Name())
		if err != nil {
			return nil, err
		}

		triggerCode := statevos.NewTrigger(int(statevos.ScopeCode))
		triggerRecipe := statevos.NewTrigger(int(statevos.ScopeRecipe))
		triggerVars := statevos.NewTrigger(int(statevos.ScopeVars))

		fingerprints := map[statevos.Trigger]statevos.Fingerprint{
			triggerCode:   fingerprintCurrentCode,
			triggerRecipe: fingerprintCurrentStep,
			triggerVars:   fingerprintCurrentVars,
		}

		decision := s.statePolicyService.Decide(stateLatest, stepDef.TriggersInt(), fingerprints)

		if decision.ShouldExecute() {
			stateCurrent := stateAgg.NewExecutionState(stepExec.Name())
			stateCurrent.SetFingerprint(triggerCode, fingerprintCurrentCode)
			stateCurrent.SetFingerprint(triggerRecipe, fingerprintCurrentStep)
			stateCurrent.SetFingerprint(triggerVars, fingerprintCurrentVars)

			err = s.processStep(req, order, stepExec, stateCurrent)
			if err != nil {
				return order, err
			}
		} else {
			order.MarkStepAsCached(stepExec.Name())
			fmt.Printf("---------------- OMITIENDO STEP: %s -----------------\n", stepExec.Name())
			fmt.Println("-------------------------------------------------------")
			continue
		}
	}

	if order.Status() == orchVos.OrderStatusFailed {
		fmt.Println("❌ La ejecución de la orden ha fallado")
	}
	if order.Status() == orchVos.OrderStatusSuccessful {
		fmt.Println("🎉 Ejecución de la orden completada con éxito.")
	}
	return order, nil
}

func (s *ExecuteOrder) processStepVariables(stepName string, stepDefinition depEnt.StepDefinition, order *orchAgg.Order) error {
	storeVars, err := s.variablesRepository.FindByStepName(stepName)
	if err != nil {
		return err
	}
	order.AddOutputsMap(storeVars)

	for _, stepVar := range stepDefinition.Variables() {
		interpolatedVar, err := s.varResolver.ResolveTemplate(stepVar.Value(), order.Outputs())
		if err != nil {
			return err
		}

		order.AddOutput(stepVar.Name(), interpolatedVar)
	}
	return nil
}

func (s *ExecuteOrder) processStep(
	req appDto.OrderRequest,
	order *orchAgg.Order,
	stepExec *orchEnt.StepRecord,
	stateCurrent *stateAgg.ExecutionState) error {

	err := s.executeStep(req, order, stepExec)
	if err != nil {
		return err
	}

	if stepExec.Status() == orchVos.StepStatusFailed {
		msm := fmt.Sprintf("❌ La ejecución del paso '%s' ha fallado", stepExec.Name())
		return errors.New(msm)
	}

	if stepExec.Status() == orchVos.StepStatusSuccessful {
		err = s.variablesRepository.Save(stepExec.Name(), order.GetOutputMap())
		if err != nil {
			return err
		}

		err = s.stateRepository.Save(stateCurrent)
		if err != nil {
			return err
		}

		if err := s.orderRepo.Save(order, req.ProjectDom.Project().Name()); err != nil {
			return err
		}

		fmt.Printf("🎉 Ejecución del paso '%s' completada con éxito.\n", stepExec.Name())
	}
	return nil
}

func (s *ExecuteOrder) executeStep(
	req appDto.OrderRequest,
	order *orchAgg.Order,
	stepExec *orchEnt.StepRecord) error {

	workdirStep, err := s.workspaceMgr.Prepare(stepExec.Name())

	if err != nil {
		return err
	}

	var_step_workdir := "step_workdir"
	var_command_workdir := "command_workdir"

	order.AddOutput(var_step_workdir, workdirStep)

	fmt.Printf("--------------- EJECUTANDO STEP: %s ---------------\n", stepExec.Name())
	fmt.Println("--------------------------------------------------------")
	for _, cmdExec := range stepExec.Commands() {
		fmt.Printf("-> Ejecutando comando: %s\n", cmdExec.Name())

		workdirCmd := cmdExec.Workdir()
		if workdirCmd != "" && workdirStep != "" {
			workdirCmd = s.cmdExecutor.CreateWorkDir(workdirStep, workdirCmd)
			order.AddOutput(var_command_workdir, workdirCmd)
		}

		interpolatedCmd, err := s.varResolver.ResolveTemplate(cmdExec.Command(), order.Outputs())
		if err != nil {
			return err
		}

		for _, templatePath := range cmdExec.TemplateFiles() {
			pathTemplateFile := s.cmdExecutor.CreateWorkDir(workdirCmd, templatePath)

			err = s.varResolver.ResolvePath(pathTemplateFile, order.Outputs())
			if err != nil {
				return err
			}
		}

		log, exitCode, err := s.cmdExecutor.Execute(req.Ctx, workdirCmd, interpolatedCmd)
		if err != nil {
			return err
		}

		order.RemoveOutput(var_command_workdir)

		err = order.FinalizeCommand(stepExec.Name(), cmdExec.Name(), interpolatedCmd, log, exitCode, s.varResolver)
		if err != nil {
			return err
		}
	}
	order.RemoveOutput(var_step_workdir)

	return nil
}

func (s *ExecuteOrder) findStepDefinition(steps []depEnt.StepDefinition, stepName string) depEnt.StepDefinition {
	for _, step := range steps {
		if step.Name() == stepName {
			return step
		}
	}
	return depEnt.StepDefinition{}
}

func (s *ExecuteOrder) getVariablesInit(req appDto.OrderRequest) ([]orchVos.Output, error) {
	environment := req.Environment.Value()

	allVarsInit := []orchVos.Output{}

	varsProject, err := s.getVariablesProject(req.ProjectDom)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, varsProject...)

	env, err := orchVos.NewOutputFromNameAndValue("environment", environment)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, env)

	projectPath, err := orchVos.NewOutputFromNameAndValue("project_workdir", req.ProjectPath)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, projectPath)

	toolName, err := orchVos.NewOutputFromNameAndValue("tool_name", "fastdeploy")
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, toolName)

	return allVarsInit, nil
}

func (s *ExecuteOrder) getVariablesProject(domModel *domAgg.DeploymentObjectModel) ([]orchVos.Output, error) {
	varsProject := []orchVos.Output{}

	projectId, err := orchVos.NewOutputFromNameAndValue("project_id", domModel.Project().IdString()[:8])
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectId)

	projectName, err := orchVos.NewOutputFromNameAndValue("project_name", domModel.Project().Name())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectName)

	projectTeam, err := orchVos.NewOutputFromNameAndValue("project_team", domModel.Project().Team())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectTeam)

	projectVersion, err := orchVos.NewOutputFromNameAndValue("project_version", domModel.Project().Version())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectVersion)

	projectRevision, err := orchVos.NewOutputFromNameAndValue("project_revision", domModel.Project().Revision())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, projectRevision)

	productId, err := orchVos.NewOutputFromNameAndValue("product_id", domModel.Product().IdString()[:8])
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, productId)

	productName, err := orchVos.NewOutputFromNameAndValue("product_name", domModel.Product().Name())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, productName)

	productTeam, err := orchVos.NewOutputFromNameAndValue("product_team", domModel.Product().Team())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, productTeam)

	productOrganization, err := orchVos.NewOutputFromNameAndValue("product_organization", domModel.Product().Organization())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, productOrganization)

	technologyType, err := orchVos.NewOutputFromNameAndValue("technology_type", domModel.Technology().TypeTechnology())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, technologyType)

	technologySolution, err := orchVos.NewOutputFromNameAndValue("technology_solution", domModel.Technology().Solution())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, technologySolution)

	technologyStack, err := orchVos.NewOutputFromNameAndValue("technology_stack", domModel.Technology().Stack())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, technologyStack)

	technologyInfrastructure, err := orchVos.NewOutputFromNameAndValue("technology_infrastructure", domModel.Technology().Infrastructure())
	if err != nil {
		return nil, err
	}
	varsProject = append(varsProject, technologyInfrastructure)

	return varsProject, nil
}
