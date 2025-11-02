package ports

import (
	"context"

	proVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/vos"
	depAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/aggregates"
	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"
)

type TemplateRepository interface {
	LoadDeployment(
		ctx context.Context,
		source proVos.Template,
		environment string) (deployment *depAgg.Deployment, err error)

	LoadEnvironments(
		ctx context.Context,
		source proVos.Template) (environments []depVos.Environment, err error)

	PathLocal(source proVos.Template) string
}
