package application

import (
	"context"
	"errors"
	"fmt"

	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	// Usaremos nuestro dominio de statedetermination consistentemente
	sdt_services "github.com/jairoprogramador/fastdeploy-core/internal/domain/statedetermination/services"
	sdt_vos "github.com/jairoprogramador/fastdeploy-core/internal/domain/statedetermination/vos"

	execAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/aggregates"
	execEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/entities"
	execSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/services"
	execVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"

	logAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	logEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/entities"

	proAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
	proPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"

	defEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/entities"
	defPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/ports"
	defVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"

	shared "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared"
)

type AppExecutionService struct {
	varResolver           execSer.ResolverService
	fingerprintService    sdt_services.FingerprintService // Corregido para usar nuestro servicio
	workspaceMgr          appPor.StepWorkspaceService
	cmdExecutor           appPor.CommandService
	variablesRepository   statePor.VariablesRepository   // TODO: Esto parece de otro dominio de estado
	fingerprintRepository statePor.FingerprintRepository // TODO: Esto parece de otro dominio de estado
	configRepository      proPor.ConfigRepository
	templateRepository    defPor.DefinitionRepository
	gitManager            appPor.GitService
	logger                appPor.LoggerService
	stateManager          *sdt_services.StateManager // Corregido para usar nuestro servicio
}

func NewAppExecutionService(
	varResolver execSer.ResolverService,
	fingerprintService sdt_services.FingerprintService,
	workspaceMgr appPor.StepWorkspaceService,
	cmdExecutor appPor.CommandService,
	variablesRepository statePor.VariablesRepository,
	fingerprintRepository statePor.FingerprintRepository,
	configRepository proPor.ConfigRepository,
	templateRepository defPor.DefinitionRepository,
	gitManager appPor.GitService,
	logger appPor.LoggerService,
	stateManager *sdt_services.StateManager,
) *AppExecutionService {
	return &AppExecutionService{
		varResolver:           varResolver,
		fingerprintService:    fingerprintService,
		workspaceMgr:          workspaceMgr,
		cmdExecutor:           cmdExecutor,
		variablesRepository:   variablesRepository,
		fingerprintRepository: fingerprintRepository,
		configRepository:      configRepository,
		templateRepository:    templateRepository,
		gitManager:            gitManager,
		logger:                logger,
		stateManager:          stateManager,
	}
}

func (s *AppExecutionService) MarkStepAsFailed(namesParams appDto.NamesParams, logger *logAgg.Logger, step *logEnt.StepRecord, stepErr error) {
	s.logger.MarkStepAsRunning(namesParams, logger, step)
	s.logger.MarkStepAsFailed(namesParams, logger, step, stepErr)
	s.logger.FinishExecution(namesParams, logger)
}

func (s *AppExecutionService) Run(request appDto.ExecutorRequest) error {
	ctx := context.Background()

	configProject, err := s.configRepository.Load(request.PathProject())
	if err != nil {
		s.logger.ShowError("load project", err)
		return nil
	}

	environments, err := s.templateRepository.LoadEnvironments(ctx, configProject.Template())
	if err != nil {
		s.logger.ShowError("load environments", err)
		return nil
	}

	environment := request.Environment()
	if !s.existsEnvironment(environments, request.Environment()) {
		if len(environments) > 0 {
			environment = environments[0].Value()
		} else {
			err := fmt.Errorf("el ambiente '%s' no se encontró en la plantilla", request.Environment())
			s.logger.ShowError("validate environment", err)
			return nil
		}
	}

	deployment, err := s.templateRepository.LoadDeployment(ctx, configProject.Template(), environment)
	if err != nil {
		s.logger.ShowError("load deployment", err)
		return nil
	}

	if !deployment.ExistsStep(request.FinalStepName()) {
		err := fmt.Errorf("el paso '%s' no se encontró en la plantilla", request.FinalStepName())
		s.logger.ShowError("validate step", err)
		return nil
	}

	stepFinalName := deployment.StepName(request.FinalStepName())

	if stepFinalName != shared.StepTest {
		revisionProject, err := s.GetCommitHash(ctx, request.PathProject())
		if err != nil {
			s.logger.ShowError(stepFinalName, err)
			return nil
		}
		project := configProject.Project().WithRevision(revisionProject)
		configProject = configProject.WithProject(project)
	}

	allVars, err := s.getVariablesInit(configProject, request.PathProject(), environment)
	if err != nil {
		s.logger.ShowError(stepFinalName, err)
		return nil
	}

	order, err := execAgg.NewExecutionRecord(
		execVos.NewExecutionID(),
		deployment,
		environment,
		stepFinalName,
		request.SkippedStepNames(),
		allVars,
	)

	if err != nil {
		s.logger.ShowError(stepFinalName, err)
		return nil
	}

	contextDataLogger := map[string]string{
		"template":    configProject.Template().URL(),
		"environment": environment,
		"target":      stepFinalName,
	}

	namesParams := appDto.NewNamesParams(configProject.Project().Name(), configProject.Template().NameTemplate())

	execLogger, err := s.logger.StartLog(namesParams, contextDataLogger, configProject.Project().Revision()) // tal vez solo hay que pasar configProject
	if err != nil {
		s.logger.ShowError(stepFinalName, err)
		return nil
	}

	fingerprintCurrentCode, err := s.fingerprintService.GenerateFromPath(request.PathProject())
	if err != nil {
		s.logger.ShowError(stepFinalName, err)
		return nil
	}

	pathTemplateLocal := s.templateRepository.PathLocal(configProject.Template())

	for _, stepRecord := range order.StepsRecord() {
		stepLog, err := s.logger.AddStep(namesParams, execLogger, stepRecord.Name())
		if err != nil {
			s.MarkStepAsFailed(namesParams, execLogger, stepLog, err)
			return nil
		}
		if stepRecord.Status() == execVos.StepStatusSkipped {
			s.logger.MarkStepAsSkipped(namesParams, execLogger, stepLog)
			continue
		}

		stepDef := s.findStepDefinition(deployment.Steps(), stepRecord.Name())
		runParams := appDto.NewRunParams(environment, stepRecord.Name())

		err = s.processStepVariables(namesParams, runParams, stepDef, order)
		if err != nil {
			s.MarkStepAsFailed(namesParams, execLogger, stepLog, err)
			return nil
		}

		fingerprintsStateStepCurrent, err := s.getFingerprintsStateStepCurrent(stepRecord.Name(),
			pathTemplateLocal, environment, fingerprintCurrentCode, order.GetOutputsMapForFingerprint())
		if err != nil {
			s.MarkStepAsFailed(namesParams, execLogger, stepLog, err)
			return nil
		}

		sdtStep, err := sdt_vos.NewStep(stepRecord.Name())
		if err != nil {
			s.MarkStepAsFailed(namesParams, execLogger, stepLog, err)
			return nil
		}

		shouldExecute, err := s.stateManager.HasStateChanged(request.PathProject(), sdtStep, fingerprintsStateStepCurrent, sdt_vos.NewCachePolicy(0))
		if err != nil {
			// Log the error but proceed with execution as a safe default
			s.logger.ShowWarning(fmt.Sprintf("Could not determine state change for step %s, proceeding with execution. Error: %v", stepRecord.Name(), err))
		}

		if shouldExecute {
			err = s.processStep(
				namesParams,
				runParams,
				order,
				ctx,
				fingerprintsStateStepCurrent,
				stepRecord,
				execLogger)

			if err != nil {
				s.logger.MarkStepAsFailed(namesParams, execLogger, stepLog, err)
				s.logger.FinishExecution(namesParams, execLogger)
				return nil
			}

			if stepRecord.Status() == execVos.StepStatusSuccessful {
				s.logger.MarkStepAsSuccessful(namesParams, execLogger, stepLog)
			} else {
				err := fmt.Errorf("estado inesperado del paso '%s': %s %s", stepRecord.Name(), stepRecord.Status().String(), err)
				s.logger.MarkStepAsFailed(namesParams, execLogger, stepLog, err)
				s.logger.FinishExecution(namesParams, execLogger)
				return nil
			}
		} else {
			order.MarkStepAsCached(stepRecord.Name())
			s.logger.MarkStepAsCached(namesParams, execLogger, stepLog, "State has not changed")
			continue
		}
	}

	s.logger.FinishExecution(namesParams, execLogger)
	return nil
}

func (s *AppExecutionService) getFingerprintsStateStepCurrent(
	stepName, templatePath, environment string,
	fingerprintCurrentCode sdt_vos.Fingerprint,
	varsMap map[string]string) (sdt_vos.CurrentStateFingerprints, error) {

	fingerprintCurrentVars, err := s.fingerprintService.GenerateFromVariables(varsMap)
	if err != nil {
		return sdt_vos.CurrentStateFingerprints{}, err
	}

	fingerprintCurrentInstruction, err := s.fingerprintService.GenerateFromStepDefinition(templatePath, appDto.NewRunParams(environment, stepName))
	if err != nil {
		return sdt_vos.CurrentStateFingerprints{}, err
	}

	// ACL: No necesitamos traducir el environment aquí porque ya lo hacemos en el StateManager.
	// Pasamos el string directamente.
	sdtEnv, err := sdt_vos.NewEnvironment(environment)
	if err != nil {
		return sdt_vos.CurrentStateFingerprints{}, err
	}

	return sdt_vos.NewCurrentStateFingerprints(
		fingerprintCurrentCode,
		fingerprintCurrentInstruction,
		fingerprintCurrentVars,
		sdtEnv,
	), nil
}

func (s *AppExecutionService) processStepVariables(namesParams appDto.NamesParams, runParams appDto.RunParams, stepDefinition defEnt.StepDefinition, order *execAgg.ExecutionRecord) error {
	storeVars, err := s.variablesRepository.FindByStep(namesParams, runParams)
	if err != nil {
		return err
	}
	order.AddOutputsMap(storeVars)

	for _, stepVar := range stepDefinition.Variables() {
		interpolatedVar, err := s.varResolver.ResolveString(stepVar.Value(), order.Outputs())
		if err != nil {
			return err
		}
		order.AddOutput(stepVar.Name(), interpolatedVar)
	}
	return nil
}

func (s *AppExecutionService) processStep(
	namesParams appDto.NamesParams,
	runParams appDto.RunParams,
	order *execAgg.ExecutionRecord,
	ctx context.Context,
	stateCurrent sdt_vos.CurrentStateFingerprints,
	stepRecord *execEnt.StepRecord,
	logger *logAgg.Logger) error {

	err := s.executeStep(namesParams, runParams, order, ctx, stepRecord, logger)
	if err != nil {
		return err
	}

	if stepRecord.Status() == execVos.StepStatusSuccessful {
		if err := s.variablesRepository.SaveByStep(namesParams, runParams, order.GetOutputsMapForSave()); err != nil {
			return err
		}

		sdtStep, err := sdt_vos.NewStep(stepRecord.Name())
		if err != nil {
			return err // Could not create a valid domain step, cannot update state
		}

		if err := s.stateManager.UpdateState(order.ProjectPath(), sdtStep, stateCurrent); err != nil {
			// Log a warning, as the step itself was successful, but state saving failed.
			s.logger.ShowWarning(fmt.Sprintf("Step %s executed successfully, but failed to update state: %v", stepRecord.Name(), err))
		}
	}
	return nil
}

func (s *AppExecutionService) executeStep(
	namesParams appDto.NamesParams,
	runParams appDto.RunParams,
	order *execAgg.ExecutionRecord,
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
			s.logger.MarkTaskAsFailed(namesParams, logger, taskLog, err, stepLog)
			return err
		}

		workdirCmd := cmdExec.Workdir()
		if workdirCmd != "" && workdirStep != "" {
			workdirCmd = s.cmdExecutor.CreateWorkDir(workdirStep, workdirCmd)
			order.AddOutput(execAgg.OutputCommWorkdirKey, workdirCmd)
		}

		interpolatedCmd, err := s.varResolver.ResolveString(cmdExec.Command(), order.Outputs())
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
		cmdOutput, exitCode, err := s.cmdExecutor.Run(ctx, workdirCmd, interpolatedCmd)
		if err != nil {
			taskLog.AddOutput(cmdOutput)
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
	stepRecord.UpdateStatus()
	return nil
}

func (s *AppExecutionService) findStepDefinition(steps []defEnt.StepDefinition, stepName string) defEnt.StepDefinition {
	for _, step := range steps {
		if step.Name() == stepName {
			return step
		}
	}
	return defEnt.StepDefinition{}
}

func (s *AppExecutionService) getVariablesInit(configProject *proAgg.MyProject, pathProject string, environment string) ([]execVos.Output, error) {

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

func (s *AppExecutionService) getVariablesConfig(config *proAgg.MyProject) ([]execVos.Output, error) {
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

func (s *AppExecutionService) GetCommitHash(ctx context.Context, pathProject string) (string, error) {
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

func (s *AppExecutionService) existsEnvironment(environments []defVos.EnvironmentDefinition, environmentValue string) bool {
	for _, env := range environments {
		if env.Value() == environmentValue {
			return true
		}
	}
	return false
}
