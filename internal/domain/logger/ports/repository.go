package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy/internal/domain/logger/aggregates"
)

type LoggerRepository interface {
	Save(namesParams dto.NamesParams, logger *aggregates.Logger) error
	Find(namesParams dto.NamesParams) (aggregates.Logger, error)
}
