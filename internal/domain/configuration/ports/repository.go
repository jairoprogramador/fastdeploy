package ports

import "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"

type Repository interface {
	Exists() (bool, error)
	Load() (entities.Configuration, error)
	Save(config entities.Configuration) error
}
