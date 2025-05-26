package main

import (
	"deploy/internal/application"
	cmd "deploy/internal/cli/command"
	"deploy/internal/cli/handler"
	"deploy/internal/domain/engine"
	"deploy/internal/domain/engine/condition"
	"deploy/internal/domain/engine/executor"
	engineModel "deploy/internal/domain/engine/model"
	"deploy/internal/domain/engine/validator"
	domainModel "deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/service"
	"deploy/internal/infrastructure/adapter"
	"deploy/internal/infrastructure/repository"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// Log message constants
const (
	logPrefix                  = "FASTDEPLOY: "
	msgInitCommon              = "Initializing common components..."
	msgInitRepo                = "Initializing repositories..."
	msgInitDomainServices      = "Initializing domain services..."
	msgInstantiatingExecutors  = "Instantiating step executors..."
	msgInstantiatingEngine     = "Instantiating engine..."
	msgRegisteringExecutors    = "Registering step executors..."
	msgInstantiatingCommands   = "Instantiating commands..."
	msgRunningCLI              = "--- Running CLI Application ---"
)

func main() {
	mainLogger := log.New(os.Stdout, logPrefix, log.LstdFlags)

	mainLogger.Println(msgInitCommon)
	appRouter := service.NewPathService()
	appLogger := logger.NewLogger(appRouter.GetFullPathLoggerFile())
	variableStore := engineModel.NewStoreEntity()

	mainLogger.Println(msgInitRepo)
	fileRepo := adapter.NewOsFileController()
	yamlRepo := adapter.NewGoPkgYamlController(fileRepo)
	configRepo := repository.NewYamlConfigRepository(yamlRepo, fileRepo, appRouter)
	projectRepo := repository.NewYamlProjectRepository(yamlRepo, fileRepo, appRouter)
	deploymentRepo := repository.NewYamlDeploymentRepository(yamlRepo, fileRepo, appRouter)

	mainLogger.Println(msgInitDomainServices)
	executorService := adapter.NewOsRunCommand()
	conditionFactory := condition.NewEvaluatorFactory()
	deploymentValidator := validator.NewValidator(appLogger)
	baseExecutor := executor.NewStepExecutor()

	// Initialize infrastructure services
	gitService := adapter.NewLocalGitRequest(executorService)
	templateService := adapter.NewTextDockerTemplate()

	// Initialize domain services
	globalConfigService := service.NewConfigService(configRepo)
	projectService := service.NewProjectService(appLogger, projectRepo, globalConfigService, appRouter)
	storeService := service.NewStoreService(appLogger, gitService, appRouter)
	deploymentService := service.NewDeploymentService(deploymentRepo, appRouter)

	// Initialize docker image service
	dockerImageService := adapter.NewLocalDockerImage(fileRepo, templateService, projectService, appRouter, variableStore)

	// Initialize docker container service with docker image dependency
	dockerService := adapter.NewLocalDockerContainer(executorService, fileRepo, templateService, dockerImageService, appRouter, variableStore, appLogger)

	mainLogger.Println(msgInstantiatingExecutors)
	commandExecutor := executor.NewCommandExecutor(appLogger, baseExecutor, variableStore, executorService, conditionFactory)
	containerExecutor := executor.NewContainerExecutor(baseExecutor, variableStore, dockerService)
	setupExecutor := executor.NewCheckExecutor(appLogger, dockerService, variableStore, appRouter)

	mainLogger.Println(msgInstantiatingEngine)
	engineInstance := engine.NewEngine(
		variableStore,
		storeService,
		appLogger,
		deploymentValidator,
	)

	mainLogger.Println(msgRegisteringExecutors)
	engineInstance.Executors[string(engineModel.Command)] = commandExecutor
	engineInstance.Executors[string(engineModel.Container)] = containerExecutor
	engineInstance.Executors[string(engineModel.Setup)] = setupExecutor

	mainLogger.Println(msgInstantiatingCommands)
	deployCmdFn := getDeployCmdFn()
	initCmd := newInitCmd(projectService)
	startCmd := newStartCmd(projectService, engineInstance, deploymentService)
	cmd.SetupCommands(deployCmdFn, initCmd, startCmd)

	mainLogger.Println(msgRunningCLI)
	cmd.Execute()
}

func getDeployCmdFn() func() *cobra.Command {
	deployHandler := handler.NewDeployHandler()
	getDeployCmdFn := func() *cobra.Command {
		return cmd.GetDeployCmd(deployHandler.Controller)
	}
	return getDeployCmdFn
}

func newInitCmd(projectService service.ProjectService) *cobra.Command {
	initAppFn := func() error {
		return application.InitApp(projectService)
	}
	initHandler := handler.NewInitHandler(initAppFn)
	return cmd.NewInitCmd(initHandler.Controller)
}

func newStartCmd(projectService service.ProjectService, engineInstance *engine.Engine, deploymentService service.DeploymentLoader) *cobra.Command {
	isInitAppFn := func() (*domainModel.ProjectEntity, error) {
		project, err := projectService.Load()
		if err != nil {
			if err == service.ErrProjectNotFound || err == service.ErrProjectNotComplete {
				return nil, err
			}
			return nil, fmt.Errorf("failed to check project initialization status: %w", err)
		}
		return project, nil
	}

	startAppFn := func(project *domainModel.ProjectEntity) error {
		return application.StartDeploy(engineInstance, deploymentService, project)
	}

	startHandler := handler.NewStartHandler(startAppFn, isInitAppFn)
	return cmd.NewStartCmd(startHandler.Controller)
}
