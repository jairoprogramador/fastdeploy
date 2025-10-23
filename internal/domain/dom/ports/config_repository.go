package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
)

type ConfigRepository interface {
	Save(config *aggregates.Config) error
	Load() (*aggregates.Config, error)
}
