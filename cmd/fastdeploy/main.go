package main

import (
	"github.com/jairoprogramador/fastdeploy/internal/application"
	cmd "github.com/jairoprogramador/fastdeploy/internal/cli/command"
	"github.com/jairoprogramador/fastdeploy/internal/cli/handler"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/condition"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/executor"
	engineModel "github.com/jairoprogramador/fastdeploy/internal/domain/engine/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/store"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
	domainModel "github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model/logger"
	"github.com/jairoprogramador/fastdeploy/internal/domain/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/repository"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// Log message constants
const (
	logPrefix                 = "FASTDEPLOY: "
	msgInitCommon             = "Initializing common components..."
	msgInitRepo               = "Initializing repositories..."
	msgInitDomainServices     = "Initializing domain services..."
	msgInstantiatingExecutors = "Instantiating step executors..."
	msgInstantiatingEngine    = "Instantiating engine..."
	msgRegisteringExecutors   = "Registering step executors..."
	msgInstantiatingCommands  = "Instantiating commands..."
	msgRunningCLI             = "--- Running CLI Application ---"
)

func main() {
	mainLogger := log.New(os.Stdout, logPrefix, log.LstdFlags)

	mainLogger.Println(msgInitCommon)
	appRouter := adapter.NewOsPathService()
	appLoggerFile := logger.NewFileLogger(appRouter.GetFullPathLoggerFile())
	variableStore := engineModel.NewStoreEntity()

	mainLogger.Println(msgInitRepo)
	fileRepo := adapter.NewOsFileController(appLoggerFile)
	yamlRepo := adapter.NewGoPkgYamlController(fileRepo, appLoggerFile)
	configRepo := repository.NewYamlConfigRepository(yamlRepo, fileRepo, appRouter, appLoggerFile)
	projectRepo := repository.NewYamlProjectRepository(yamlRepo, fileRepo, appRouter, appLoggerFile)
	deploymentRepo := repository.NewYamlDeploymentRepository(yamlRepo, fileRepo, appRouter, appLoggerFile)

	mainLogger.Println(msgInitDomainServices)
	executorService := adapter.NewOsRunCommand(appLoggerFile)
	conditionFactory := condition.NewEvaluatorFactory()
	deploymentValidator := validator.NewValidator()
	baseExecutor := executor.NewStepExecutor()

	// Initialize infrastructure services
	gitService := adapter.NewLocalGitRequest(executorService, appLoggerFile)
	templateService := adapter.NewTextDockerTemplate(appLoggerFile)

	// Initialize domain services
	deploymentService := service.NewDeploymentService(deploymentRepo)
	configService := service.NewConfigService(configRepo)

	storeService := store.NewStoreService(gitService, appRouter)

	mainLogger.Println(msgInstantiatingEngine)
	engineInstance := engine.NewEngine(
		variableStore,
		storeService,
		deploymentValidator,
	)
	projectService := service.NewProjectService(projectRepo, deploymentService, engineInstance, configService, appRouter)

	// Initialize docker image service
	dockerImageService := adapter.NewLocalDockerImage(fileRepo, templateService, projectService, appRouter, variableStore)

	// Initialize docker container service with docker image dependency
	dockerService := adapter.NewLocalDockerContainer(executorService, fileRepo, templateService, dockerImageService, appRouter, variableStore, appLoggerFile)

	mainLogger.Println(msgInstantiatingExecutors)
	commandExecutor := executor.NewCommandExecutor(baseExecutor, variableStore, executorService, conditionFactory)
	containerExecutor := executor.NewContainerExecutor(baseExecutor, variableStore, dockerService)
	setupExecutor := executor.NewCheckExecutor(dockerService, variableStore, appRouter)

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
	initAppFn := func() domainModel.DomainResultEntity {
		return application.InitApp(projectService)
	}
	initHandler := handler.NewInitHandler(initAppFn)
	return cmd.NewInitCmd(initHandler.Controller)
}

func newStartCmd(projectService service.ProjectService, engineInstance *engine.Engine, deploymentService service.DeploymentService) *cobra.Command {
	
	startAppFn := func() domainModel.DomainResultEntity {
		return application.StartDeploy(projectService)
	}

	startHandler := handler.NewStartHandler(startAppFn)
	return cmd.NewStartCmd(startHandler.Controller)
}
