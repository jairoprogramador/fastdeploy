package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
)

type Writer interface {
	Write(entities.Project) error
}
