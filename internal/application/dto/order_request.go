package dto

import (
	"context"

	depAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/aggregates"
	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"
	domAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
)

type OrderRequest struct {
	Ctx              context.Context
	Environment      depVos.Environment
	FinalStep        string
	Template         *depAgg.DeploymentTemplate
	TemplatePath     string
	RepositoryName   string
	ProjectDom       *domAgg.Config
	ProjectPath      string
	SkippedStepNames map[string]struct{}
}
