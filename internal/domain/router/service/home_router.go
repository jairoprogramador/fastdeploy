package service

import "github.com/jairoprogramador/fastdeploy/internal/domain/router/port"


const FASTDEPLOY_HOME_VARIABLE = "FASTDEPLOY_HOME"
const FASTDEPLOY_DIR = ".fastdeploy"

type HomeRouterService interface {
	GetPathWorkdir() (string, error)
	GetPathFastDeploy() (string, error)
}

type HomeRouterServiceImpl struct {
	homeRouter port.HomeRouter
}

func NewHomeRouterService(homeRouter port.HomeRouter) HomeRouterService {
	return &HomeRouterServiceImpl{homeRouter: homeRouter}
}

func (h *HomeRouterServiceImpl) GetPathFastDeploy() (string, error) {
	environmentVariable := h.homeRouter.GetEnvironmentVariable(FASTDEPLOY_HOME_VARIABLE)

	if environmentVariable == "" {

		currentUserDir, err := h.homeRouter.GetCurrentUserDir()
		if err != nil {
			return "", err
		}

		homeRoute := h.homeRouter.BuildRoute(currentUserDir, FASTDEPLOY_DIR)
		return homeRoute, nil
	}

	return environmentVariable, nil
}

func (p *HomeRouterServiceImpl) GetPathWorkdir() (string, error) {
	return p.homeRouter.GetPathWorkdir()
}
