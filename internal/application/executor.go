package application

import (
	"context"
	"errors"
	"fmt"

	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	stateAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	statePor "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
	stateSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/services"
	statevos "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"

	execAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/executor/aggregates"
	execEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/executor/entities"
	execSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/executor/services"
	execVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/executor/vos"

	logAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"

	proAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
	proPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"

	temEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/entities"
	temPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/ports"
	temVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"

	shared "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared"
)

type AppExecutor struct {
	varResolver           execSer.TemplateResolver
	fingerprintService    stateSer.FingerprintService
	workspaceMgr          appPor.StepWorkspace
	cmdExecutor           appPor.CommandExecutor
	variablesRepository   statePor.VariablesRepository
	fingerprintRepository statePor.FingerprintRepository
	statePolicyService    stateSer.FingerprintPolicyService
	configRepository      proPor.ConfigRepository
	templateRepository    temPor.TemplateRepository
	gitManager            appPor.GitManager
	logger                appPor.Logger
}

func NewAppExecutor(
	varResolver execSer.TemplateResolver,
	fingerprintService stateSer.FingerprintService,
	workspaceMgr appPor.StepWorkspace,
	cmdExecutor appPor.CommandExecutor,
	variablesRepository statePor.VariablesRepository,
	fingerprintRepository statePor.FingerprintRepository,
	statePolicyService stateSer.FingerprintPolicyService,
	configRepository proPor.ConfigRepository,
	templateRepository temPor.TemplateRepository,
	gitManager appPor.GitManager,
	logger appPor.Logger,
) *AppExecutor {
	return &AppExecutor{
		varResolver:           varResolver,
		fingerprintService:    fingerprintService,
		workspaceMgr:          workspaceMgr,
		cmdExecutor:           cmdExecutor,
		variablesRepository:   variablesRepository,
		fingerprintRepository: fingerprintRepository,
		statePolicyService:    statePolicyService,
		configRepository:      configRepository,
		templateRepository:    templateRepository,
		gitManager:            gitManager,
		logger:                logger,
	}
}

func (s *AppExecutor) existsEnvironment(environments []temVos.Environment, environmentValue string) bool {
	for _, env := range environments {
		if env.Value() == environmentValue {
			return true
		}
	}
	return false
}

func (s *AppExecutor) Run(request appDto.ExecutorRequest) error {

	ctx := context.Background()

	configProject, err := s.configRepository.Load(request.PathProject())
	if err != nil {
		return err
	}

	environments, err := s.templateRepository.LoadEnvironments(ctx, configProject.Template())
	if err != nil {
		return err
	}

	environment := request.Environment()
	if !s.existsEnvironment(environments, request.Environment()) {
		if len(environments) >= 0 {
			environment = environments[0].Value()
		} else {
			return fmt.Errorf("el ambiente '%s' no se encontró en la plantilla", request.Environment())
		}
	}

	deployment, err := s.templateRepository.LoadDeployment(ctx, configProject.Template(), environment)
	if err != nil {
		return err
	}

	if !deployment.ExistsStep(request.FinalStepName()) {
		return fmt.Errorf("el paso '%s' no se encontró en la plantilla", request.FinalStepName())
	}

	stepFinalName := deployment.StepName(request.FinalStepName())

	if stepFinalName != shared.StepTest {
		revisionProject, err := s.GetCommitHash(ctx, request.PathProject())
		if err != nil {
			return err
		}
		configProject.SetProjectRevision(revisionProject)
	}

	allVars, err := s.getVariablesInit(configProject, request.PathProject(), environment)
	if err != nil {
		return err
	}

	fingerprintCurrentCode, err := s.fingerprintService.GenerateFromPath(request.PathProject())
	if err != nil {
		return err
	}

	order, err := execAgg.NewOrder(
		execVos.NewOrderID(),
		deployment,
		environment,
		stepFinalName,
		request.SkippedStepNames(),
		allVars,
	)

	if err != nil {
		return err
	}

	contextDataLogger := map[string]string{
		"template":    configProject.Template().URL(),
		"environment": environment,
		"target":      stepFinalName,
	}

	namesParams := appDto.NewNamesParams(configProject.Project().Name(), configProject.Template().NameTemplate())

	execLogger, err := s.logger.Start(namesParams, contextDataLogger, configProject.Project().Revision()) // tal vez solo hay que pasar configProject
	if err != nil {
		return err
	}

	pathTemplateLocal := s.templateRepository.PathLocal(configProject.Template())

	for _, stepRecord := range order.StepsRecord() {
		stepLog, err := s.logger.AddStep(namesParams, execLogger, stepRecord.Name())
		if err != nil {
			return err
		}

		if stepRecord.Status() == execVos.StepStatusSkipped {
			s.logger.MarkStepAsSkipped(namesParams, execLogger, stepLog)
			continue
		}

		stepDef := s.findStepDefinition(deployment.Steps(), stepRecord.Name())
		runParams := appDto.NewRunParams(environment, stepRecord.Name())

		err = s.processStepVariables(namesParams, runParams, stepDef, order)
		if err != nil {
			s.logger.MarkStepAsFailed(namesParams, execLogger, stepLog, err)
			s.logger.FinishExecution(namesParams, execLogger)
			return err
		}

		fingerprintsStateStepCurrent, err := s.getFingerprintsStateStepCurrent(stepRecord.Name(),
			pathTemplateLocal, environment, fingerprintCurrentCode, order.GetOutputsMapForFingerprint())
		if err != nil {
			return err
		}

		fingerprintsStateStepLatest, err := s.getFingerprintsStateStepLatest(namesParams, runParams)
		if err != nil {
			return err
		}

		decision := s.statePolicyService.Decide(fingerprintsStateStepLatest, stepDef.TriggersInt(), fingerprintsStateStepCurrent)

		if decision.ShouldExecute() {
			err = s.processStep(
				namesParams,
				runParams,
				order,
				ctx,
				fingerprintsStateStepCurrent,
				stepRecord,
				execLogger)

			if err != nil {
				s.logger.FinishExecution(namesParams, execLogger)
				return err
			}
		} else {
			order.MarkStepAsCached(stepRecord.Name())
			s.logger.MarkStepAsCached(namesParams, execLogger, stepLog, decision.Reason())
			continue
		}
	}

	s.logger.FinishExecution(namesParams, execLogger)
	if order.Status() == execVos.OrderStatusFailed {
		return errors.New("❌ la ejecución de la orden ha fallado")
	}

	if order.Status() == execVos.OrderStatusSuccessful {
		fingerprintsStateCodeCurrent := stateAgg.NewFingerprintState(statevos.ScopeCode.String())
		fingerprintsStateCodeCurrent.SetFingerprint(statevos.ScopeCode, fingerprintCurrentCode)

		err = s.fingerprintRepository.SaveCode(namesParams, fingerprintsStateCodeCurrent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AppExecutor) getFingerprintsStateStepCurrent(
	stepName, templatePath, environment string,
	fingerprintCurrentCode statevos.Fingerprint,
	varsMap map[string]string) (*stateAgg.FingerprintState, error) {

	fingerprintCurrentVars, err := s.fingerprintService.GenerateFromStepVariables(varsMap)
	if err != nil {
		return &stateAgg.FingerprintState{}, err
	}

	fingerprintCurrentRecipe, err := s.fingerprintService.GenerateFromStepDefinition(templatePath, appDto.NewRunParams(environment, stepName))
	if err != nil {
		return &stateAgg.FingerprintState{}, err
	}

	fingerprintsStateStepCurrent := stateAgg.NewFingerprintState(stepName)
	fingerprintsStateStepCurrent.SetFingerprint(statevos.ScopeCode, fingerprintCurrentCode)
	fingerprintsStateStepCurrent.SetFingerprint(statevos.ScopeRecipe, fingerprintCurrentRecipe)
	fingerprintsStateStepCurrent.SetFingerprint(statevos.ScopeVars, fingerprintCurrentVars)

	return fingerprintsStateStepCurrent, nil
}

func (s *AppExecutor) getFingerprintsStateStepLatest(namesParams appDto.NamesParams, runParams appDto.RunParams) (*stateAgg.FingerprintState, error) {

	fingerprintsStateStepLatest, err := s.fingerprintRepository.FindStep(namesParams, runParams)
	if err != nil {
		return nil, err
	}
	fingerprintsStateCodeLatest, err := s.fingerprintRepository.FindCode(namesParams)
	if err != nil {
		return nil, err
	}
	fingerprintsCodeLatest, ok := fingerprintsStateCodeLatest.GetFingerprint(statevos.ScopeCode)
	if ok {
		fingerprintsStateStepLatest.SetFingerprint(statevos.ScopeCode, fingerprintsCodeLatest)
	}

	return fingerprintsStateStepLatest, nil
}

func (s *AppExecutor) processStepVariables(namesParams appDto.NamesParams, runParams appDto.RunParams, stepDefinition temEnt.StepDefinition, order *execAgg.Order) error {
	storeVars, err := s.variablesRepository.FindByStepName(namesParams, runParams)
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

func (s *AppExecutor) processStep(
	namesParams appDto.NamesParams,
	runParams appDto.RunParams,
	order *execAgg.Order,
	ctx context.Context,
	stateCurrent *stateAgg.FingerprintState,
	stepRecord *execEnt.StepRecord,
	logger *logAgg.Logger) error {

	stepLog, _ := logger.GetStep(stepRecord.Name())

	err := s.executeStep(namesParams, runParams, order, ctx, stepRecord, logger)
	if err != nil {
		return err
	}

	if stepRecord.Status() == execVos.StepStatusSuccessful {
		err = s.variablesRepository.Save(namesParams, runParams, order.GetOutputsMapForSave())
		if err != nil {
			s.logger.MarkStepAsFailed(namesParams, logger, stepLog, err)
			return err
		}

		err = s.fingerprintRepository.SaveStep(namesParams, runParams, stateCurrent)
		if err != nil {
			s.logger.MarkStepAsFailed(namesParams, logger, stepLog, err)
			return err
		}

		s.logger.MarkStepAsSuccessful(namesParams, logger, stepLog)
	}
	return nil
}

func (s *AppExecutor) executeStep(
	namesParams appDto.NamesParams,
	runParams appDto.RunParams,
	order *execAgg.Order,
	ctx context.Context,
	stepRecord *execEnt.StepRecord,
	logger *logAgg.Logger) error {

	workdirStep, err := s.workspaceMgr.Prepare(namesParams, runParams)

	if err != nil {
		return err
	}

	order.AddOutput(execAgg.OutputStepWorkdirKey, workdirStep)

	stepLog, _ := logger.GetStep(stepRecord.Name())
	s.logger.MarkStepAsRunning(namesParams, logger, stepLog)

	for _, cmdExec := range stepRecord.Commands() {
		taskLog, err := s.logger.AddTaskToStep(namesParams, logger, stepRecord.Name(), cmdExec.Name())
		if err != nil {
			return err
		}

		workdirCmd := cmdExec.Workdir()
		if workdirCmd != "" && workdirStep != "" {
			workdirCmd = s.cmdExecutor.CreateWorkDir(workdirStep, workdirCmd)
			order.AddOutput(execAgg.OutputCommWorkdirKey, workdirCmd)
		}

		interpolatedCmd, err := s.varResolver.ResolveTemplate(cmdExec.Command(), order.Outputs())
		if err != nil {
			s.logger.MarkTaskAsFailed(namesParams, logger, taskLog, err, stepLog)
			return err
		}
		s.logger.SetTaskCommand(namesParams, logger, taskLog, interpolatedCmd)

		for _, templatePath := range cmdExec.TemplateFiles() {
			pathTemplateFile := s.cmdExecutor.CreateWorkDir(workdirCmd, templatePath)

			err = s.varResolver.ResolvePath(pathTemplateFile, order.Outputs())
			if err != nil {
				s.logger.MarkTaskAsFailed(namesParams, logger, taskLog, err, stepLog)
				return err
			}
		}

		s.logger.MarkTaskAsRunning(namesParams, logger, taskLog, stepLog)
		cmdOutput, exitCode, err := s.cmdExecutor.Execute(ctx, workdirCmd, interpolatedCmd)
		if err != nil {
			s.logger.MarkTaskAsFailed(namesParams, logger, taskLog, err, stepLog)
			return err
		}

		s.logger.AddOutputToTask(namesParams, logger, taskLog, cmdOutput)
		err = order.FinalizeCommand(stepRecord.Name(), cmdExec.Name(), interpolatedCmd, cmdOutput, exitCode, s.varResolver)
		if err != nil {
			s.logger.MarkTaskAsFailed(namesParams, logger, taskLog, err, stepLog)
			return err
		}
		if exitCode != 0 {
			err = fmt.Errorf("el comando finalizó con código de salida %d", exitCode)
			s.logger.MarkTaskAsFailed(namesParams, logger, taskLog, err, stepLog)
			return err
		}
		s.logger.MarkTaskAsSuccessful(namesParams, logger, taskLog, stepLog)
	}
	return nil
}

func (s *AppExecutor) findStepDefinition(steps []temEnt.StepDefinition, stepName string) temEnt.StepDefinition {
	for _, step := range steps {
		if step.Name() == stepName {
			return step
		}
	}
	return temEnt.StepDefinition{}
}

func (s *AppExecutor) getVariablesInit(configProject *proAgg.Config, pathProject string, environment string) ([]execVos.Output, error) {

	allVarsInit := []execVos.Output{}

	varsProject, err := s.getVariablesConfig(configProject)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, varsProject...)

	env, err := execVos.NewOutputFromNameAndValue("environment", environment)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, env)

	projectPath, err := execVos.NewOutputFromNameAndValue("project_workdir", pathProject)
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, projectPath)

	toolName, err := execVos.NewOutputFromNameAndValue("tool_name", "fastdeploy")
	if err != nil {
		return nil, err
	}
	allVarsInit = append(allVarsInit, toolName)

	return allVarsInit, nil
}

func (s *AppExecutor) getVariablesConfig(config *proAgg.Config) ([]execVos.Output, error) {
	varsDeployment := []execVos.Output{}

	projectId, err := execVos.NewOutputFromNameAndValue("project_id", config.Project().IdString()[:8])
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectId)

	projectName, err := execVos.NewOutputFromNameAndValue("project_name", config.Project().Name())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectName)

	projectTeam, err := execVos.NewOutputFromNameAndValue("project_team", config.Project().Team())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectTeam)

	projectVersion, err := execVos.NewOutputFromNameAndValue("project_version", config.Project().Version())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectVersion)

	projectRevision, err := execVos.NewOutputFromNameAndValue("project_revision", config.Project().Revision())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectRevision)

	projectOrganization, err := execVos.NewOutputFromNameAndValue("project_organization", config.Project().Organization())
	if err != nil {
		return nil, err
	}
	varsDeployment = append(varsDeployment, projectOrganization)

	return varsDeployment, nil
}

func (s *AppExecutor) GetCommitHash(ctx context.Context, pathProject string) (string, error) {
	isGit, err := s.gitManager.IsGit(pathProject)
	if err != nil {
		return "", err
	}

	if !isGit {
		return "", errors.New("el projecto no esta configurado como repositorio git, ejecute 'git init' primero")
	}

	existsChanges, err := s.gitManager.ExistChanges(ctx, pathProject)
	if err != nil {
		return "", err
	}

	if existsChanges {
		return "", errors.New("hay cambios en el proyecto, ejecute 'git commit' primero")
	}

	return s.gitManager.GetCommitHash(ctx, pathProject)
}
