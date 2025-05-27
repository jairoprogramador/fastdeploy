package main

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	cmdCli "github.com/jairoprogramador/fastdeploy/internal/cli/command"
	"github.com/jairoprogramador/fastdeploy/internal/cli/handler"
	serviceConfig "github.com/jairoprogramador/fastdeploy/internal/domain/config/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	serviceDeploy "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/condition"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/executor"
	serviceEngine "github.com/jairoprogramador/fastdeploy/internal/domain/engine/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
	serviceProject "github.com/jairoprogramador/fastdeploy/internal/domain/project/service"
	cmdAdapter "github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/command"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/docker"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/git"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/path"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/template"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/yaml"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/repository"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"github.com/spf13/cobra"
	"log"
	"os"
)

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
	appPath := path.NewPathAdapter()
	appLoggerFile := logger.NewFileLogger(appPath.GetFullPathLoggerFile())
	store := model.NewStoreEntity()

	// Initialize infrastructure services
	commandAdapter := cmdAdapter.NewCommandAdapter(appLoggerFile)
	gitAdapter := git.NewGitAdapter(commandAdapter, appLoggerFile)
	templateAdapter := template.NewTemplateAdapter(appLoggerFile)
	fileAdapter := file.NewFileAdapter(appLoggerFile)
	yamlAdapter := yaml.NewYamlAdapter(fileAdapter, appLoggerFile)

	mainLogger.Println(msgInitRepo)
	configRepository := repository.NewConfigRepository(yamlAdapter, fileAdapter, appPath, appLoggerFile)
	projectRepository := repository.NewProjectRepository(yamlAdapter, fileAdapter, appPath, appLoggerFile)
	deployRepository := repository.NewDeploymentRepository(yamlAdapter, fileAdapter, appPath, appLoggerFile)

	mainLogger.Println(msgInitDomainServices)
	evaluatorFactory := condition.NewEvaluatorFactory()
	validator := validator.NewValidator()
	baseExecutor := executor.NewBaseExecutor()
	deploymentService := serviceDeploy.NewDeploymentService(deployRepository)
	configService := serviceConfig.NewConfigService(configRepository)
	storeService := serviceEngine.NewStoreService(gitAdapter, appPath)

	mainLogger.Println(msgInstantiatingEngine)
	engineInstance := engine.NewEngine(
		store,
		storeService,
		validator,
	)
	projectService := serviceProject.NewProjectService(projectRepository, deploymentService, engineInstance, configService)
	imageAdapter := docker.NewImageAdapter(fileAdapter, templateAdapter, projectRepository, appPath, store)
	containerAdapter := docker.NewContainerAdapter(commandAdapter, fileAdapter, templateAdapter, imageAdapter, appPath, store, appLoggerFile)

	mainLogger.Println(msgInstantiatingExecutors)
	commandExecutor := executor.NewCommandExecutor(baseExecutor, commandAdapter, evaluatorFactory)
	containerExecutor := executor.NewContainerExecutor(baseExecutor, containerAdapter)
	checkExecutor := executor.NewCheckExecutor(baseExecutor, containerAdapter, store)

	mainLogger.Println(msgRegisteringExecutors)
	engineInstance.AddExecutor(model.Command, commandExecutor)
	engineInstance.AddExecutor(model.Container, containerExecutor)
	engineInstance.AddExecutor(model.Check, checkExecutor)

	mainLogger.Println(msgInstantiatingCommands)
	deployCmdFn := getDeployCmdFn()
	initCmd := newInitCmd(projectService)
	startCmd := newStartCmd(projectService)
	cmdCli.SetupCommands(deployCmdFn, initCmd, startCmd)

	mainLogger.Println(msgRunningCLI)
	cmdCli.Execute()
}

func getDeployCmdFn() func() *cobra.Command {
	deployHandler := handler.NewDeployHandler()
	getDeployCmdFn := func() *cobra.Command {
		return cmdCli.GetDeployCmd(deployHandler.Controller)
	}
	return getDeployCmdFn
}

func newInitCmd(projectService serviceProject.ProjectService) *cobra.Command {
	initAppFn := func() result.DomainResult {
		return project.InitApp(projectService)
	}
	initHandler := handler.NewInitHandler(initAppFn)
	return cmdCli.NewInitCmd(initHandler.Controller)
}

func newStartCmd(projectService serviceProject.ProjectService) *cobra.Command {

	startAppFn := func() result.DomainResult {
		return project.StartDeploy(projectService)
	}

	startHandler := handler.NewStartHandler(startAppFn)
	return cmdCli.NewStartCmd(startHandler.Controller)
}
