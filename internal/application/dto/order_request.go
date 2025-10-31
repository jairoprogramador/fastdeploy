package dto

import (
	"context"

	depAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/aggregates"
	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"
	domAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
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
