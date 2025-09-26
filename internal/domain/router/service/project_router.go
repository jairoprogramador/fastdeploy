package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/router/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/router/values"
)

const PROJECT_HOME_DIR = "projects"

type ProjectRouterService interface {
	GetPathProject() string
	GetPathEnvironment() string
	GetPathStep(step string) string
	BuildPath(paths ...string) string
}

type ProjectRouterServiceImpl struct {
	projectRouter port.ProjectRouter
	parameter     values.Parameter
}

func NewProjectRouterService(
	projectRouter port.ProjectRouter,
	parameter values.Parameter) ProjectRouterService {

	return &ProjectRouterServiceImpl{projectRouter: projectRouter, parameter: parameter}
}

func (p *ProjectRouterServiceImpl) BuildPath(paths ...string) string {
	return p.projectRouter.BuildRoute(paths...)
}

func (p *ProjectRouterServiceImpl) GetPathProject() string {
	return p.projectRouter.BuildRoute(
		p.parameter.GetHomeDir(),
		PROJECT_HOME_DIR,
		p.parameter.GetProjectName(),
	)
}

func (p *ProjectRouterServiceImpl) GetPathEnvironment() string {
	pathProject := p.GetPathProject()

	return p.projectRouter.BuildRoute(
		pathProject,
		p.parameter.GetEnvironment(),
	)
}

func (p *ProjectRouterServiceImpl) GetPathStep(step string) string {
	pathEnvironment := p.GetPathEnvironment()

	return p.projectRouter.BuildRoute(
		pathEnvironment,
		step,
	)
}
