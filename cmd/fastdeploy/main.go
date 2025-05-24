package main

import (
	"deploy/internal/application"
	"deploy/internal/cli/command"
	"deploy/internal/cli/handler"
	"deploy/internal/domain/engine"
	"deploy/internal/domain/engine/condition"
	executor2 "deploy/internal/domain/engine/executor"
	"deploy/internal/domain/engine/validator"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/service"
	"deploy/internal/domain/service/router"
	infraService "deploy/internal/infrastructure/adapter"
	"deploy/internal/infrastructure/repository"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func main() {
	mainLogger := log.New(os.Stdout, "FASTDEPLOY: ", log.LstdFlags)

	mainLogger.Println("Initializing common components...")
	appRouter := router.NewRouter()
	appLogger := logger.NewLogger(appRouter.GetFullPathLoggerFile())
	variableStore := model.NewVariableStore()

	mainLogger.Println("Initializing repositories...")
	fileRepo := infraService.NewFileRepositoryImpl()
	yamlRepo := infraService.NewYamlRepositoryImpl(fileRepo)
	configRepo := repository.NewYamlConfigRepository(yamlRepo, fileRepo, appRouter)
	projectRepo := repository.NewYamlProjectRepository(yamlRepo, fileRepo, appRouter)
	deploymentRepo := repository.NewYamlDeploymentRepository(yamlRepo, fileRepo, appRouter)
	//containerRepo := repository.NewContainerRepositoryImpl(fileRepo)

	mainLogger.Println("Initializing domain services...")
	executorInfraService := infraService.NewExecutorService()
	conditionFactory := condition.NewConditionFactory()
	deploymentValidator := validator.NewDeploymentValidator(appLogger)
	baseExecutor := executor2.NewBaseExecutor()

	// Initialize infrastructure services
	gitInfraService := infraService.NewLocalGitCommand(executorInfraService)
	templateService := infraService.NewTextDockerTemplate()

	// Initialize domain services
	globalConfigService := service.NewConfigService(configRepo)
	projectService := service.NewProjectService(appLogger, projectRepo, globalConfigService, appRouter)
	storeService := service.NewStoreService(appLogger, gitInfraService, appRouter)
	deploymentService := service.NewDeploymentService(deploymentRepo, appRouter)

	// Initialize docker image service
	dockerImageService := infraService.NewLocalDockerImage(executorInfraService, fileRepo, templateService, projectService, appRouter, variableStore)

	// Initialize docker container service with docker image dependency
	dockerInfraService := infraService.NewLocalDockerContainer(executorInfraService, fileRepo, templateService, dockerImageService, appRouter, variableStore, appLogger)

	mainLogger.Println("Instantiating step executors...")
	commandExecutor := executor2.NewCommandExecutor(appLogger, baseExecutor, variableStore, executorInfraService, conditionFactory)
	containerExecutor := executor2.NewContainerExecutor(baseExecutor, variableStore, dockerInfraService)
	setupExecutor := executor2.NewSetupExecutor(appLogger, dockerInfraService, variableStore, appRouter)

	mainLogger.Println("Instantiating engine...")
	engineInstance := engine.NewEngine(
		variableStore,
		storeService,
		appLogger,
		deploymentValidator,
	)

	mainLogger.Println("Registering step executors...")
	engineInstance.Executors[validator.TypeCommand] = commandExecutor
	engineInstance.Executors[validator.TypeContainer] = containerExecutor
	engineInstance.Executors[validator.TypeSetup] = setupExecutor

	mainLogger.Println("Instantiating commands...")
	deployCmdFn := getDeployCmdFn()
	initCmd := newInitCmd(projectService)
	startCmd := newStartCmd(projectService, engineInstance, deploymentService)
	addCmd := newAddCmd()
	cmd.SetupCommands(deployCmdFn, initCmd, startCmd, addCmd)

	mainLogger.Println("--- Running CLI Application ---")
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
	isInitAppFn := func() (*model.ProjectEntity, error) {
		project, err := projectService.Load()
		if err != nil {
			if err == service.ErrProjectNotFound || err == service.ErrProjectNotComplete {
				return nil, err
			}
			return nil, fmt.Errorf("failed to check project initialization status: %w", err)
		}
		return project, nil
	}

	startAppFn := func(project *model.ProjectEntity) error {
		return application.StartDeploy(engineInstance, deploymentService, project)
	}

	startHandler := handler.NewStartHandler(startAppFn, isInitAppFn)
	return cmd.NewStartCmd(startHandler.Controller)
}

func newAddCmd() *cobra.Command {
	addSupportHandler := newAddSupportHandler()
	addDependencyHandler := newAddDependencyHandler()
	return cmd.NewAddCmd(
		addSupportHandler.ControllerSonarQube,
		addSupportHandler.ControllerFortify,
		addDependencyHandler.Controller,
	)
}

func newAddSupportHandler() *handler.AddSupportHandler {
	addSonarQubeAppFn := func() (string, error) {
		return "", fmt.Errorf("no implementado")
	}
	addFortifyAppFn := func() (string, error) {
		return "", fmt.Errorf("no implementado")
	}

	return handler.NewAddSupportHandler(addSonarQubeAppFn, addFortifyAppFn)
}

func newAddDependencyHandler() *handler.AddDependencyHandler {
	addProjectDependencyAppFn := func(name string, version string) (string, error) {
		return "", fmt.Errorf("no implementado")
	}
	return handler.NewAddDependencyHandler(addProjectDependencyAppFn)
}
