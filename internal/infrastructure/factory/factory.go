package factory

import (
	"fmt"
	"os"
	"path/filepath"

	applic "github.com/jairoprogramador/fastdeploy-core/internal/application"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/services"
	aggregates "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/aggregates"
	servicesExec "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/services"
	dStateServices "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/services"
	dVersionServices "github.com/jairoprogramador/fastdeploy-core/internal/domain/versioning/services"
	iDefin "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/definition"
	iExecu "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/execution"
	iState "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/state"
	iGitRepo "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/versioning"
	iLogSer "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/logger/service"
	iLogRep "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/logger/repository"
	iProje "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project"

	"github.com/spf13/viper"
)

type ServiceFactory interface {
	BuildExecutionOrchestrator() (*applic.ExecutionOrchestrator, error)
	BuildLogService() *applic.LoggerService
	PathAppProject() string
}

type Factory struct {
	pathAppProject    string
	pathAppFastdeploy string
}

func NewFactory() (ServiceFactory, error) {
	fastdeployHome := getFastdeployHome()

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error al obtener el directorio de trabajo: %w", err)
	}

	return &Factory{
		pathAppFastdeploy: fastdeployHome,
		pathAppProject:    workingDir,
	}, nil
}

func (f *Factory) PathAppProject() string {
	return f.pathAppProject
}

func (f *Factory) BuildLogService() *applic.LoggerService {
	consolePresenter := iLogSer.NewConsolePresenterService()
	loggerRepository := iLogRep.NewFileLoggerRepository("")
	configRepository := iProje.NewYAMLProjectRepository()

	return applic.NewLoggerService(loggerRepository, configRepository, consolePresenter)
}

func (f *Factory) BuildExecutionOrchestrator() (*applic.ExecutionOrchestrator, error) {
	// Infrastructure Layer
	commandRunner := iExecu.NewShellCommandRunner()
	fileSystem := iExecu.NewOSFileSystem()
	gitClonerTemplate := iProje.NewGitClonerTemplate()
	gitRepository := iGitRepo.NewGoGitRepository()
	definitionReader := iDefin.NewYamlDefinitionReader()
	projectRepository := iProje.NewYAMLProjectRepository()
	fingerprintService := iState.NewSha256FingerprintService()
	stateRepository := iState.NewGobStateRepository()
	copyWorkdir := iExecu.NewCopyWorkdir()

	// Domain & Application Services
	projectService := applic.NewProjectService(projectRepository)
	workspaceService := applic.NewWorkspaceService()
	versionCalculator := dVersionServices.NewVersionCalculator(gitRepository)
	planBuilder := services.NewPlanBuilder(definitionReader)
	stateManager := dStateServices.NewStateManager(stateRepository)
	interpolator := servicesExec.NewInterpolator()
	fileProcessor := servicesExec.NewFileProcessor(fileSystem, interpolator)
	outputExtractor := servicesExec.NewOutputExtractor()
	commandExecutor := aggregates.NewCommandExecutor(commandRunner, fileProcessor, interpolator, outputExtractor)
	stepExecutor := aggregates.NewStepExecutor(commandExecutor)

	orchestrator := applic.NewExecutionOrchestrator(
		f.pathAppProject,
		f.pathAppFastdeploy,
		projectService,
		workspaceService,
		gitClonerTemplate,
		versionCalculator,
		planBuilder,
		fingerprintService,
		stateManager,
		stepExecutor,
		copyWorkdir,
	)
	return orchestrator, nil
}

func getFastdeployHome() string {
	viper.SetEnvPrefix("FASTDEPLOY")
	viper.AutomaticEnv()

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error al obtener el directorio home:", err)
		os.Exit(1)
	}

	defaultHome := filepath.Join(userHomeDir, ".fastdeploy")
	fastdeployHome := viper.GetString("HOME")
	if fastdeployHome == "" {
		fastdeployHome = defaultHome
	}
	return fastdeployHome
}
