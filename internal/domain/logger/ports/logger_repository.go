package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/vos"
)

type LoggerRepository interface {
	Save(logger *aggregates.Logger) error
	FindByID(loggerID vos.LoggerID) (*aggregates.Logger, error)
}
