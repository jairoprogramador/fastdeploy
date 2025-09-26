package router

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/router/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/shared"
)

type RepositoryRouterSO struct {}

func NewRepositoryRouterSO() port.RepositoryRouter {
	return &RepositoryRouterSO{}
}

func (rrs *RepositoryRouterSO) BuildRoute(paths ...string) string {
	return shared.GetPath(paths...)
}