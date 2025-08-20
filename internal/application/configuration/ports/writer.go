package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
)

type Writer interface {
	Write(entities.Configuration) error
}
