package main

import (
	"log"
	"os"
	"deploy/internal/cli/command"
	"deploy/internal/application"
	"deploy/internal/domain/condition"
	"deploy/internal/domain/engine"
	"deploy/internal/domain/executor"
	"deploy/internal/domain/model"
	"deploy/internal/domain/router"
	"deploy/internal/domain/service"
	"deploy/internal/domain/validator"
	"deploy/internal/infrastructure/repository"
	"deploy/internal/cli/handler"
	"fmt"
	infraService "deploy/internal/infrastructure/service"
	"github.com/spf13/cobra"
)

func main() {
	mainLogger := log.New(os.Stdout, "FASTDEPLOY: ", log.LstdFlags)

	mainLogger.Println("Initializing common components...")
	appRouter := router.NewRouter()
	appLogStore := model.NewLogStore("FastDeploy")
	variableStore := model.NewVariableStore()

	mainLogger.Println("Initializing repositories...")
	fileRepo := repository.NewFileRepositoryImpl()
	yamlRepo := repository.NewYamlRepositoryImpl(fileRepo)
	containerRepo := repository.NewContainerRepositoryImpl(fileRepo)

	mainLogger.Println("Initializing domain services...")
	executorInfraService := infraService.NewExecutorService(appLogStore)
	conditionFactory := condition.NewConditionFactory()
	deploymentValidator := validator.NewDeploymentValidator()
	baseExecutor := executor.NewBaseExecutor()
	gitInfraService := infraService.NewGitServiceImpl(executorInfraService)
	dockerInfraService := infraService.NewDockerServiceImpl(executorInfraService)
	globalConfigService := service.NewGlobalConfigService(yamlRepo, fileRepo, appRouter)
	projectService := service.NewProjectService(yamlRepo, globalConfigService, fileRepo, appRouter, appLogStore)
	storeService := service.NewStoreService(gitInfraService, appRouter)
	deploymentService := service.NewDeploymentService(yamlRepo, fileRepo, appRouter)

	mainLogger.Println("Instantiating step executors...")
	commandExecutor := executor.NewCommandExecutor(baseExecutor, variableStore,executorInfraService, conditionFactory)
	containerExecutor := executor.NewContainerExecutor(baseExecutor, variableStore, dockerInfraService, containerRepo, fileRepo, appRouter)
	setupExecutor := executor.NewSetupExecutor(dockerInfraService, variableStore, appRouter, appLogStore)
	
	mainLogger.Println("Instantiating engine...")
	engineInstance := engine.NewEngine(
		variableStore,
		storeService,
		appLogStore,
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

func getDeployCmdFn() (func() *cobra.Command){
	deployHandler := handler.NewDeployHandler()
	getDeployCmdFn := func() *cobra.Command {
		return cmd.GetDeployCmd(deployHandler.Controller)
	}
	return getDeployCmdFn
}

func newInitCmd(projectService service.ProjectServiceInterface) *cobra.Command{
	initAppFn := func(appLogger *model.LogStore) *model.LogStore {
		return application.InitApp(projectService, appLogger)
	}
	initHandler := handler.NewInitHandler(initAppFn)
	return cmd.NewInitCmd(initHandler.Controller)
}

func newStartCmd(projectService service.ProjectServiceInterface, engineInstance *engine.Engine, deploymentService service.DeploymentServiceInterface) *cobra.Command{
	isInitAppFn := func() (*model.Project, error) {
		project, err := projectService.Load()
		if err != nil {
			if err == service.ErrProjectNotFound || err == service.ErrProjectNotComplete {
				return nil, err
			}
			return nil, fmt.Errorf("failed to check project initialization status: %w", err)
		}
		return project, nil
	}

	startAppFn := func(project *model.Project) *model.LogStore {
		return application.StartDeploy(engineInstance, deploymentService, project)
	}

	startHandler := handler.NewStartHandler(startAppFn, isInitAppFn)
	return cmd.NewStartCmd(startHandler.Controller)
}

func newAddCmd() *cobra.Command{
	addSupportHandler := newAddSupportHandler()
	addDependencyHandler := newAddDependencyHandler()
	return cmd.NewAddCmd(
		addSupportHandler.ControllerSonarQube,
		addSupportHandler.ControllerFortify,
		addDependencyHandler.Controller,
	)
}

func newAddSupportHandler() *handler.AddSupportHandler{
	addSonarQubeAppFn := func() (string, error) {
		return "", fmt.Errorf("no implementado")
	}
	addFortifyAppFn := func() (string, error) {
		return "", fmt.Errorf("no implementado")
	}

	return handler.NewAddSupportHandler(addSonarQubeAppFn, addFortifyAppFn)
}

func newAddDependencyHandler() *handler.AddDependencyHandler{
	addProjectDependencyAppFn := func(name string, version string) (string, error) {
		return "", fmt.Errorf("no implementado")
	}
	return handler.NewAddDependencyHandler(addProjectDependencyAppFn)
}
