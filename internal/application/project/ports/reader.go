package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
)

type Reader interface {
	ExistsFile() (bool, error)
	Read() (entities.Project, error)
}