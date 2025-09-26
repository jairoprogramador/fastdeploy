package router

import (
	"os"
	"os/user"

	"github.com/jairoprogramador/fastdeploy/internal/domain/router/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/shared"
)

type HomeRouterSO struct {}

func NewHomeRouterSO() port.HomeRouter {
	return &HomeRouterSO{}
}

func (fsr *HomeRouterSO) GetEnvironmentVariable(nameVariable string) string {
	return os.Getenv(nameVariable)
}

func (fsr *HomeRouterSO) GetCurrentUserDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}

func (fsr *HomeRouterSO) BuildRoute(paths ...string) string {
	return shared.GetPath(paths...)
}

func (prs *HomeRouterSO) GetPathWorkdir() (string, error) {
	return os.Getwd()
}
