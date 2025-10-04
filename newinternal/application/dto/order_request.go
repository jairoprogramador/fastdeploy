package dto

import (
	"context"

	deploymentaggregates "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	domaggregates "github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/aggregates"
)

type OrderRequest struct {
	Ctx              context.Context
	Environment      vos.Environment
	FinalStep        string
	Template         *deploymentaggregates.DeploymentTemplate
	TemplatePath     string
	RepositoryName   string
	ProjectDom       *domaggregates.DeploymentObjectModel
	ProjectPath      string
	SkippedStepNames map[string]struct{}
}
