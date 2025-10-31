package factory

import (
	"path/filepath"
	"github.com/spf13/viper"
	"fmt"
	"os"
	applic "github.com/jairoprogramador/fastdeploy-core/internal/application"
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"
	iLogge "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/logger"
	iProje "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project"
)

type ServiceFactory interface {
	BuildLogService() appPor.Logger
}

type Factory struct {
}

func NewFactory() ServiceFactory {
	return &Factory{}
}

func (f *Factory) BuildLogService() appPor.Logger {
	fastdeployHome := getFastdeployHome()
	statePath := filepath.Join(fastdeployHome, "state")

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	domRepository := iProje.NewDomYAMLRepository(workingDir)
	loadDOMService := applic.NewLoadConfigService(domRepository)
	config, err := loadDOMService.Load()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	consolePresenter := iLogge.NewConsolePresenter()
	loggerRepository, err := iLogge.NewFileLoggerRepository(
		statePath,
		config.Project().Name(),
		config.Template().NameTemplate())
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	return applic.NewAppLogger(loggerRepository, consolePresenter)
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