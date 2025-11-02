package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
)

type ConfigRepository interface {
	Save(config *aggregates.Config, pathProject string) error
	Load(pathProject string) (*aggregates.Config, error)
}
