package factory

import (
	"fmt"
	"os"
	"path/filepath"

	applic "github.com/jairoprogramador/fastdeploy-core/internal/application"
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	staSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/services"

	iLogge "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/logger"
	iProje "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project"

	iExecu "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/executor"

	iStaRep "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/state/repository"
	iStaSer "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/state/services"

	iAppli "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/application"

	iTempl "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/template"

	"github.com/spf13/viper"
)

type ServiceFactory interface {
	BuildLogService() appPor.Logger
	BuildExecutorService() *applic.AppExecutor
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

func (f *Factory) BuildLogService() appPor.Logger {
	consolePresenter := iLogge.NewConsolePresenter()
	loggerRepository := iLogge.NewFileLoggerRepository(f.pathStateRoot)
	configRepository := iProje.NewYamlConfigRepository()

	return applic.NewAppLogger(loggerRepository, configRepository, consolePresenter)
}

func (f *Factory) BuildExecutorService() *applic.AppExecutor {
	varResolver := iExecu.NewGoTemplateResolver()
	fingerprintService := iStaSer.NewFingerprintService()
	fileManager := iAppli.NewFileStepWorkspace(f.pathProjectsRoot, f.pathRepositoriesRoot)
	cmdExecutor := iAppli.NewExecutor()
	varsRepository := iStaRep.NewVarsRepository(f.pathStateRoot)
	stateRepository := iStaRep.NewFileFingerprintRepository(f.pathStateRoot)
	statePolicyService := staSer.NewFingerprintPolicyService()
	configRepository := iProje.NewYamlConfigRepository()
	templateRepository := iTempl.NewTemplateRepository(f.pathRepositoriesRoot, cmdExecutor)
	gitManager := iAppli.NewGitManager(cmdExecutor)

	loggerRepository := iLogge.NewFileLoggerRepository(f.pathStateRoot)
	consolePresenter := iLogge.NewConsolePresenter()
	logger := applic.NewAppLogger(loggerRepository, configRepository, consolePresenter)

	return applic.NewAppExecutor(
		varResolver,
		fingerprintService,
		fileManager,
		cmdExecutor,
		varsRepository,
		stateRepository,
		statePolicyService,
		configRepository,
		templateRepository,
		gitManager,
		logger,
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
