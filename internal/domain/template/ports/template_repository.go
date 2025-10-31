package ports

import (
	"context"

	shaVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
	depAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/aggregates"
)

type TemplateRepository interface {
	Load(ctx context.Context, source shaVos.TemplateSource) (template *depAgg.DeploymentTemplate, repositoryLocalPath string, err error)
}
