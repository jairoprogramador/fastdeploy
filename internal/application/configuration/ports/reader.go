package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
)

type Reader interface {
	Read() (entities.Configuration, error)
}