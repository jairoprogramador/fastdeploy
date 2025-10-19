package ports

import (
	"context"

	sharedVos "github.com/jairoprogramador/fastdeploy/internal/domain/shared/vos"

	depAgg "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/aggregates"
)

type TemplateRepository interface {
	Load(ctx context.Context, source sharedVos.TemplateSource) (template *depAgg.DeploymentTemplate, repositoryLocalPath string, err error)
}
