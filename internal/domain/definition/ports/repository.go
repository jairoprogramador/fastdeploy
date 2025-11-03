package ports

import (
	"context"

	proVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/vos"
	defAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/aggregates"
	defVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
)

type DefinitionRepository interface {
	LoadDeployment(
		ctx context.Context,
		source proVos.Template,
		environment string) (deployment *defAgg.Deployment, err error)

	LoadEnvironments(
		ctx context.Context,
		source proVos.Template) (environments []defVos.EnvironmentDefinition, err error)

	PathLocal(source proVos.Template) string
}
