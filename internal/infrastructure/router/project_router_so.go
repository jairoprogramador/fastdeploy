package router

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/router/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/shared"
)

type ProjectRouterSO struct {}

func NewProjectRouterSO() port.ProjectRouter {
	return &ProjectRouterSO{}
}

func (prs *ProjectRouterSO) BuildRoute(paths ...string) string {
	return shared.GetPath(paths...)
}