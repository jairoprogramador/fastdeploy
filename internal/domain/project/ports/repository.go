package ports

import "github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"

type Repository interface {
	Exists() (bool, error)
	Load() (entities.Project, error)
	Save(project entities.Project) error
}
