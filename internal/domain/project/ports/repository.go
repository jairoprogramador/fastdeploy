package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
)

type ConfigRepository interface {
	Save(config *aggregates.MyProject, pathProject string) error
	Load(pathProject string) (*aggregates.MyProject, error)
}
