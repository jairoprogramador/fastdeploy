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
)

const (
	msgInitCommon           = "Common components initialized"
	msgInitInfrastructure   = "Initialized infrastructure services"
	msgInitDomainServices   = "Domain services Initialized"
	msgInitExecutors        = "Initialized executors"
	msgRegisteredExecutors  = "Registered executors"
	msgInstantiatedCommands = "Instantiated commands"
)

func main() {
	appPath := path.NewPathAdapter()
	appLoggerFile := logger.NewFileLogger(appPath.GetFullPathLoggerFile())
	store := model.NewStoreEntity()

	appLoggerFile.Info(msgInitCommon)

	commandAdapter := cmdAdapter.NewCommandAdapter(appLoggerFile)
	gitAdapter := git.NewGitAdapter(commandAdapter, appLoggerFile)
	templateAdapter := template.NewTemplateAdapter(appLoggerFile)
	fileAdapter := file.NewFileAdapter(appLoggerFile)
	yamlAdapter := yaml.NewYamlAdapter(fileAdapter, appLoggerFile)
	appLoggerFile.Info(msgInitInfrastructure)

	configRepository := repository.NewConfigRepository(yamlAdapter, fileAdapter, appPath, appLoggerFile)
	projectRepository := repository.NewProjectRepository(yamlAdapter, fileAdapter, appPath, appLoggerFile)
	deployRepository := repository.NewDeploymentRepository(yamlAdapter, fileAdapter, appPath, appLoggerFile)

	evaluatorFactory := condition.NewEvaluatorFactory()
	validator := validator.NewValidator()
	baseExecutor := executor.NewBaseExecutor()
	configService := serviceConfig.NewConfigService(configRepository)
	storeService := serviceEngine.NewStoreService(gitAdapter, appPath)

	engineInstance := engine.NewEngine(
		store,
		validator,
	)
	imageAdapter := docker.NewImageAdapter(fileAdapter, templateAdapter, projectRepository, appPath, store)
	containerAdapter := docker.NewContainerAdapter(commandAdapter, fileAdapter, templateAdapter, imageAdapter, appPath, store, appLoggerFile)
	deploymentService := serviceDeploy.NewDeploymentService(deployRepository, containerAdapter, store)
	projectService := serviceProject.NewProjectService(projectRepository, deploymentService, engineInstance, configService, containerAdapter, store, storeService)

	appLoggerFile.Info(msgInitDomainServices)

	commandExecutor := executor.NewCommandExecutor(baseExecutor, commandAdapter, evaluatorFactory)
	containerExecutor := executor.NewContainerExecutor(baseExecutor, containerAdapter, store)
	checkExecutor := executor.NewCheckExecutor(baseExecutor, containerAdapter)
	appLoggerFile.Info(msgInitExecutors)

	engineInstance.AddExecutor(model.Command, commandExecutor)
	engineInstance.AddExecutor(model.Container, containerExecutor)
	engineInstance.AddExecutor(model.Check, checkExecutor)
	appLoggerFile.Info(msgRegisteredExecutors)

	deployCmdFn := getDeployCmdFn()
	initCmd := newInitCmd(projectService, appLoggerFile)
	startCmd := newStartCmd(projectService, appLoggerFile)
	cmdCli.SetupCommands(deployCmdFn, initCmd, startCmd)
	appLoggerFile.Info(msgInstantiatedCommands)

	cmdCli.Execute()
}

func getDeployCmdFn() func() *cobra.Command {
	deployHandler := handler.NewDeployHandler()
	getDeployCmdFn := func() *cobra.Command {
		return cmdCli.GetDeployCmd(deployHandler.Controller)
	}
	return getDeployCmdFn
}

func newInitCmd(projectService serviceProject.ProjectService, fileLogger *logger.FileLogger) *cobra.Command {
	initAppFn := func() result.DomainResult {
		return project.InitApp(projectService, fileLogger)
	}
	initHandler := handler.NewInitHandler(initAppFn)
	return cmdCli.NewInitCmd(initHandler.Controller, fileLogger)
}

func newStartCmd(projectService serviceProject.ProjectService, fileLogger *logger.FileLogger) *cobra.Command {

	startAppFn := func() result.DomainResult {
		return project.StartDeploy(projectService, fileLogger)
	}

	startHandler := handler.NewStartHandler(startAppFn)
	return cmdCli.NewStartCmd(startHandler.Controller, fileLogger)
}
