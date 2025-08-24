package port

import "github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"

type Repository interface {
	Exists() (bool, error)
	Load() (entity.Project, error)
	Save(project entity.Project) error
}
