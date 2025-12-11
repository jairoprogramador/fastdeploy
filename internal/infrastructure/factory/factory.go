package factory

import (
	"fmt"
	"os"
	"path/filepath"

	applic "github.com/jairoprogramador/fastdeploy-core/internal/application"
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	dStatePorts "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
	dStateServices "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/services"

	iLogRep "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/logger/repository"
	iLogSer "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/logger/service"

	iProje "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project"

	iExecu "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/execution"

	iStaRep "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/state/repository"
	iStaSer "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/state/services"

	iAppli "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/application"

	iDefin "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/definition"

	"github.com/spf13/viper"
)

type ServiceFactory interface {
	BuildLogService() appPor.LoggerService
	BuildExecutorService() *applic.AppExecutionService
	PathAppProject() string
}

type Factory struct {
	pathRepositoriesRoot string
	pathProjectsRoot     string
	pathStateRoot        string
	pathAppProject       string
}

func NewFactory() (ServiceFactory, error) {
	fastdeployHome := getFastdeployHome()

	pathRepositoriesRoot := filepath.Join(fastdeployHome, "repositories")
	pathProjectsRoot := filepath.Join(fastdeployHome, "projects")
	pathStateRoot := filepath.Join(fastdeployHome, "state")
	workingDir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	return &Factory{
		pathRepositoriesRoot: pathRepositoriesRoot,
		pathProjectsRoot:     pathProjectsRoot,
		pathStateRoot:        pathStateRoot,
		pathAppProject:       workingDir,
	}, nil
}

func (f *Factory) PathAppProject() string {
	return f.pathAppProject
}

func (f *Factory) BuildLogService() appPor.LoggerService {
	consolePresenter := iLogSer.NewConsolePresenterService()
	loggerRepository := iLogRep.NewFileLoggerRepository(f.pathStateRoot)
	configRepository := iProje.NewYamlConfigRepository()

	return applic.NewAppLoggerService(loggerRepository, configRepository, consolePresenter)
}

func (f *Factory) BuildExecutorService() *applic.AppExecutionService {
	varResolver := iExecu.NewResolverService()
	fingerprintGenerator := iStaSer.NewShaFingerprintGenerator()
	fileManager := iAppli.NewFileStepWorkspaceService(f.pathProjectsRoot, f.pathRepositoriesRoot)
	cmdExecutor := iAppli.NewExecCommandService()
	// Mantenemos este repo de variables por ahora, ya que parece ser para otro propósito
	varsRepository := iStaRep.NewFileVarsRepository(f.pathStateRoot)
	configRepository := iProje.NewYamlConfigRepository()
	templateRepository := iDefin.NewYamlTemplateRepository(f.pathRepositoriesRoot, cmdExecutor)
	gitManager := iAppli.NewGitLocalService(cmdExecutor)

	loggerRepository := iLogRep.NewFileLoggerRepository(f.pathStateRoot)
	consolePresenter := iLogSer.NewConsolePresenterService()
	logger := applic.NewAppLoggerService(loggerRepository, configRepository, consolePresenter)

	// Construimos nuestro dominio de `statedetermination` (ahora `state`)
	stateRepository := iStaRep.NewFileStateRepository(f.pathProjectsRoot)
	stateManager := dStateServices.NewStateManager(stateRepository)

	return applic.NewAppExecutionService(
		varResolver,
		fingerprintGenerator,
		fileManager,
		cmdExecutor,
		varsRepository,
		configRepository,
		templateRepository,
		gitManager,
		logger,
		stateManager,
	)
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
