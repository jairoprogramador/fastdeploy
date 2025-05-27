package main

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	cmd "github.com/jairoprogramador/fastdeploy/internal/cli/command"
	"github.com/jairoprogramador/fastdeploy/internal/cli/handler"
	service3 "github.com/jairoprogramador/fastdeploy/internal/domain/config/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/condition"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/executor"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/store"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
	//domainModel "github.com/jairoprogramador/fastdeploy/internal/domain/model"
	service2 "github.com/jairoprogramador/fastdeploy/internal/domain/project/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/command"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/docker"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/git"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/path"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/template"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/yaml"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/repository"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	resultEntity "github.com/jairoprogramador/fastdeploy/pkg/common/result"
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
	appRouter := path.NewOsPathService()
	appLoggerFile := logger.NewFileLogger(appRouter.GetFullPathLoggerFile())
	variableStore := entity.NewStoreEntity()

	mainLogger.Println(msgInitRepo)
	fileRepo := file.NewOsFileController(appLoggerFile)
	yamlRepo := yaml.NewGoPkgYamlController(fileRepo, appLoggerFile)
	configRepo := repository.NewYamlConfigRepository(yamlRepo, fileRepo, appRouter, appLoggerFile)
	projectRepo := repository.NewYamlProjectRepository(yamlRepo, fileRepo, appRouter, appLoggerFile)
	deploymentRepo := repository.NewYamlDeploymentRepository(yamlRepo, fileRepo, appRouter, appLoggerFile)

	mainLogger.Println(msgInitDomainServices)
	executorService := command.NewOsRunCommand(appLoggerFile)
	conditionFactory := condition.NewEvaluatorFactory()
	deploymentValidator := validator.NewValidator()
	baseExecutor := executor.NewStepExecutor()

	// Initialize infrastructure services
	gitService := git.NewLocalGitRequest(executorService, appLoggerFile)
	templateService := template.NewTextDockerTemplate(appLoggerFile)

	// Initialize domain services
	deploymentService := service.NewDeploymentService(deploymentRepo)
	configService := service3.NewConfigService(configRepo)

	storeService := store.NewStoreService(gitService, appRouter)

	mainLogger.Println(msgInstantiatingEngine)
	engineInstance := engine.NewEngine(
		variableStore,
		storeService,
		deploymentValidator,
	)
	projectService := service2.NewProjectService(projectRepo, deploymentService, engineInstance, configService, appRouter)

	// Initialize docker image service
	dockerImageService := docker.NewLocalDockerImage(fileRepo, templateService, projectService, appRouter, variableStore)

	// Initialize docker container service with docker image dependency
	dockerService := docker.NewLocalDockerContainer(executorService, fileRepo, templateService, dockerImageService, appRouter, variableStore, appLoggerFile)

	mainLogger.Println(msgInstantiatingExecutors)
	commandExecutor := executor.NewCommandExecutor(baseExecutor, variableStore, executorService, conditionFactory)
	containerExecutor := executor.NewContainerExecutor(baseExecutor, variableStore, dockerService)
	setupExecutor := executor.NewCheckExecutor(dockerService, variableStore, appRouter)

	mainLogger.Println(msgRegisteringExecutors)
	engineInstance.Executors[string(entity.Command)] = commandExecutor
	engineInstance.Executors[string(entity.Container)] = containerExecutor
	engineInstance.Executors[string(entity.Setup)] = setupExecutor

	mainLogger.Println(msgInstantiatingCommands)
	deployCmdFn := getDeployCmdFn()
	initCmd := newInitCmd(projectService)
	startCmd := newStartCmd(projectService)
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

func newInitCmd(projectService service2.ProjectService) *cobra.Command {
	initAppFn := func() resultEntity.DomainResult {
		return project.InitApp(projectService)
	}
	initHandler := handler.NewInitHandler(initAppFn)
	return cmd.NewInitCmd(initHandler.Controller)
}

func newStartCmd(projectService service2.ProjectService) *cobra.Command {

	startAppFn := func() resultEntity.DomainResult {
		return project.StartDeploy(projectService)
	}

	startHandler := handler.NewStartHandler(startAppFn)
	return cmd.NewStartCmd(startHandler.Controller)
}
