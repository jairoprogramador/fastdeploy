package factory

import (
	"fmt"
	"os"
	"path/filepath"

	applic "github.com/jairoprogramador/fastdeploy/internal/application"
	defServ "github.com/jairoprogramador/fastdeploy/internal/domain/definition/services"
	exeServ "github.com/jairoprogramador/fastdeploy/internal/domain/execution/services"
	staServ "github.com/jairoprogramador/fastdeploy/internal/domain/state/services"
	verServ "github.com/jairoprogramador/fastdeploy/internal/domain/versioning/services"
	iDefini "github.com/jairoprogramador/fastdeploy/internal/infrastructure/definition"
	iExecut "github.com/jairoprogramador/fastdeploy/internal/infrastructure/execution"
	iState "github.com/jairoprogramador/fastdeploy/internal/infrastructure/state"
	iVersi "github.com/jairoprogramador/fastdeploy/internal/infrastructure/versioning"
	iLgSer "github.com/jairoprogramador/fastdeploy/internal/infrastructure/logger/service"
	iLgRep "github.com/jairoprogramador/fastdeploy/internal/infrastructure/logger/repository"
	iProje "github.com/jairoprogramador/fastdeploy/internal/infrastructure/project"

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
	consolePresenter := iLgSer.NewConsolePresenterService()
	loggerRepository := iLgRep.NewFileLoggerRepository("")
	configRepository := iProje.NewYAMLProjectRepository()

	return applic.NewLoggerService(loggerRepository, configRepository, consolePresenter)
}

func (f *Factory) BuildExecutionOrchestrator() (*applic.ExecutionOrchestrator, error) {
	// Infrastructure Layer
	commandRunner := iExecut.NewShellCommandRunner()
	fileSystem := iExecut.NewOSFileSystem()
	gitClonerTemplate := iProje.NewGitClonerTemplate()
	gitRepository := iVersi.NewGoGitRepository()
	definitionReader := iDefini.NewYamlDefinitionReader()
	projectRepository := iProje.NewYAMLProjectRepository()
	fingerprintService := iState.NewSha256FingerprintService()
	stateRepository := iState.NewGobStateRepository()
	copyWorkdir := iExecut.NewCopyWorkdir()
	varsRepository := iExecut.NewGobVarsRepository()

	// Domain & Application Services
	projectService := applic.NewProjectService(projectRepository)
	workspaceService := applic.NewWorkspaceService()
	versionCalculator := verServ.NewVersionCalculator(gitRepository)
	planBuilder := defServ.NewPlanBuilder(definitionReader)
	stateManager := staServ.NewStateManager(stateRepository)
	interpolator := exeServ.NewInterpolator()
	fileProcessor := exeServ.NewFileProcessor(fileSystem, interpolator)
	outputExtractor := exeServ.NewOutputExtractor()
	commandExecutor := exeServ.NewCommandExecutor(commandRunner, fileProcessor, interpolator, outputExtractor)
	variableResolver := exeServ.NewVariableResolver(interpolator)
	stepExecutor := exeServ.NewStepExecutor(commandExecutor, variableResolver)

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
		varsRepository,
		gitRepository,
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
