package application

import (
	"errors"
	"fmt"

	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	stateAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	statePor "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
	stateSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/services"
	statevos "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"

	orchAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/aggregates"
	orchEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/entities"
	orchPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/ports"
	orchSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/services"
	orchVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/vos"

	depEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/entities"
	domAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
)

type ExecuteOrder struct {
	orderRepo             orchPor.OrderRepository
	varResolver           orchSer.TemplateResolver
	fingerprintService    stateSer.FingerprintService
	workspaceMgr          appPor.WorkspaceManager
	cmdExecutor           appPor.CommandExecutor
	variablesRepository   statePor.VariablesRepository
	fingerprintRepository statePor.FingerprintRepository
	statePolicyService    stateSer.FingerprintPolicyService
}

func NewExecuteOrder(
	orderRepo orchPor.OrderRepository,
	varResolver orchSer.TemplateResolver,
	fingerprintService stateSer.FingerprintService,
	workspaceMgr appPor.WorkspaceManager,
	cmdExecutor appPor.CommandExecutor,
	variablesRepository statePor.VariablesRepository,
	fingerprintRepository statePor.FingerprintRepository,
	statePolicyService stateSer.FingerprintPolicyService,
) *ExecuteOrder {
	return &ExecuteOrder{
		orderRepo:             orderRepo,
		varResolver:           varResolver,
		fingerprintService:    fingerprintService,
		workspaceMgr:          workspaceMgr,
		cmdExecutor:           cmdExecutor,
		variablesRepository:   variablesRepository,
		fingerprintRepository: fingerprintRepository,
		statePolicyService:    statePolicyService,
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

		fingerprintsStateStepCurrent, err := s.getFingerprintsStateStepCurrent(stepExec.Name(),
			req.TemplatePath, fingerprintCurrentCode, order.GetOutputsMapForFingerprint())
		if err != nil {
			return nil, err
		}

		fingerprintsStateStepLatest, err := s.getFingerprintsStateStepLatest(stepExec.Name())
		if err != nil {
			return nil, err
		}

		decision := s.statePolicyService.Decide(fingerprintsStateStepLatest, stepDef.TriggersInt(), fingerprintsStateStepCurrent)

		if decision.ShouldExecute() {
			fmt.Printf("--------------- EJECUTANDO STEP: %s ---------------\n", stepExec.Name())
			fmt.Printf("------------- %s -------------\n", decision.Reason())

			err = s.processStep(req, order, stepExec, fingerprintsStateStepCurrent)
			if err != nil {
				return order, err
			}
		} else {
			order.MarkStepAsCached(stepExec.Name())
			fmt.Printf("---------------- OMITIENDO STEP: %s -----------------\n", stepExec.Name())
			fmt.Printf("------------- %s -------------\n", decision.Reason())
			continue
		}
	}

	if order.Status() == orchVos.OrderStatusFailed {
		fmt.Println("❌ La ejecución de la orden ha fallado")
	}
	if order.Status() == orchVos.OrderStatusSuccessful {
		fingerprintsStateCodeCurrent := stateAgg.NewFingerprintState(statevos.ScopeCode.String())
		fingerprintsStateCodeCurrent.SetFingerprint(statevos.ScopeCode, fingerprintCurrentCode)

		err = s.fingerprintRepository.SaveCode(fingerprintsStateCodeCurrent)
		if err != nil {
			return order, err
		}

		fmt.Println("🎉 Ejecución de la orden completada con éxito.")
	}
	return order, nil
}

func (s *ExecuteOrder) getFingerprintsStateStepCurrent(
	stepName, templatePath string,
	fingerprintCurrentCode statevos.Fingerprint,
	varsMap map[string]string) (*stateAgg.FingerprintState, error) {

	fingerprintCurrentVars, err := s.fingerprintService.GenerateFromStepVariables(varsMap)
	if err != nil {
		return &stateAgg.FingerprintState{}, err
	}

	fingerprintCurrentRecipe, err := s.fingerprintService.GenerateFromStepDefinition(templatePath, stepName)
	if err != nil {
		return &stateAgg.FingerprintState{}, err
	}

	fingerprintsStateStepCurrent := stateAgg.NewFingerprintState(stepName)
	fingerprintsStateStepCurrent.SetFingerprint(statevos.ScopeCode, fingerprintCurrentCode)
	fingerprintsStateStepCurrent.SetFingerprint(statevos.ScopeRecipe, fingerprintCurrentRecipe)
	fingerprintsStateStepCurrent.SetFingerprint(statevos.ScopeVars, fingerprintCurrentVars)

	return fingerprintsStateStepCurrent, nil
}

func (s *ExecuteOrder) getFingerprintsStateStepLatest(stepName string) (*stateAgg.FingerprintState, error) {

	fingerprintsStateStepLatest, err := s.fingerprintRepository.FindStep(stepName)
	if err != nil {
		return nil, err
	}
	fingerprintsStateCodeLatest, err := s.fingerprintRepository.FindCode()
	if err != nil {
		return nil, err
	}
	fingerprintsCodeLatest, ok := fingerprintsStateCodeLatest.GetFingerprint(statevos.ScopeCode)
	if ok {
		fingerprintsStateStepLatest.SetFingerprint(statevos.ScopeCode, fingerprintsCodeLatest)
	}

	return fingerprintsStateStepLatest, nil
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
	stateCurrent *stateAgg.FingerprintState) error {

	err := s.executeStep(req, order, stepExec)
	if err != nil {
		return err
	}

	if stepExec.Status() == orchVos.StepStatusFailed {
		msm := fmt.Sprintf("❌ La ejecución del paso '%s' ha fallado", stepExec.Name())
		return errors.New(msm)
	}

	if stepExec.Status() == orchVos.StepStatusSuccessful {
		err = s.variablesRepository.Save(stepExec.Name(), order.GetOutputsMapForSave())
		if err != nil {
			return err
		}

		err = s.fingerprintRepository.SaveStep(stateCurrent)
		if err != nil {
			return err
		}

		if err := s.orderRepo.Save(order); err != nil {
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

	order.AddOutput(orchAgg.OutputStepWorkdirKey, workdirStep)

	for _, cmdExec := range stepExec.Commands() {
		fmt.Printf("-> Ejecutando comando: %s\n", cmdExec.Name())

		workdirCmd := cmdExec.Workdir()
		if workdirCmd != "" && workdirStep != "" {
			workdirCmd = s.cmdExecutor.CreateWorkDir(workdirStep, workdirCmd)
			order.AddOutput(orchAgg.OutputCommWorkdirKey, workdirCmd)
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

		err = order.FinalizeCommand(stepExec.Name(), cmdExec.Name(), interpolatedCmd, log, exitCode, s.varResolver)
		if err != nil {
			return err
		}
		order.RemoveOutput(orchAgg.OutputCommWorkdirKey)
	}

	order.RemoveOutput(orchAgg.OutputStepWorkdirKey)
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

	varsProject, err := s.getVariablesConfig(req.ProjectDom)
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

func (s *ExecuteOrder) getVariablesConfig(config *domAgg.Config) ([]orchVos.Output, error) {
	varsDeployment := []orchVos.Output{}

	projectId, err := orchVos.NewOutputFromNameAndValue("project_id", config.Project().IdString()[:8])
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectId)

	projectName, err := orchVos.NewOutputFromNameAndValue("project_name", config.Project().Name())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectName)

	projectTeam, err := orchVos.NewOutputFromNameAndValue("project_team", config.Project().Team())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectTeam)

	projectVersion, err := orchVos.NewOutputFromNameAndValue("project_version", config.Project().Version())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectVersion)

	projectRevision, err := orchVos.NewOutputFromNameAndValue("project_revision", config.Project().Revision())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectRevision)

	return varsDeployment, nil
}
