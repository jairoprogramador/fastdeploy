package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/router/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/router/values"
)

const REPOSITORY_HOME_DIR = "repositories"
const ENVIRONMENT_REPOSITORY_DIR = "environments"
const STEP_REPOSITORY_DIR = "steps"

type RepositoryRouterService interface {
	GetPathRepository() string
	GetPathEnvironment() string
	GetPathStep(step string) string
	BuildPath(paths ...string) string
}

type RepositoryRouterServiceImpl struct {
	repositoryRouter port.RepositoryRouter
	parameter     values.Parameter
}

func NewRepositoryRouterService(
	repositoryRouter port.RepositoryRouter,
	parameter values.Parameter) RepositoryRouterService {

	return &RepositoryRouterServiceImpl{
		repositoryRouter: repositoryRouter,
		parameter: parameter,
	}
}

func (p *RepositoryRouterServiceImpl) BuildPath(paths ...string) string {
	return p.repositoryRouter.BuildRoute(paths...)
}

func (p *RepositoryRouterServiceImpl) GetPathRepository() string {
	pathRepository := p.repositoryRouter.BuildRoute(
		p.parameter.GetHomeDir(),
		REPOSITORY_HOME_DIR,
		p.parameter.GetRepositoryName(),
	)
	return pathRepository
}

func (p *RepositoryRouterServiceImpl) GetPathEnvironment() string {
	pathRepository := p.GetPathRepository()

	pathEnvironment := p.repositoryRouter.BuildRoute(
		pathRepository,
		ENVIRONMENT_REPOSITORY_DIR,
	)
	return pathEnvironment
}

func (p *RepositoryRouterServiceImpl) GetPathStep(step string) string {
	pathRepository := p.GetPathRepository()

	if p.parameter.GetStackName() != "" {
		return p.repositoryRouter.BuildRoute(
			pathRepository,
			STEP_REPOSITORY_DIR,
			p.parameter.GetStackName(),
			step,
		)
	} else {
		return p.repositoryRouter.BuildRoute(
			pathRepository,
			STEP_REPOSITORY_DIR,
			step,
		)
	}
}
