package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
)

type LoggerRepository interface {
	Save(logger *aggregates.Logger) error
	Find() (aggregates.Logger, error)
}
