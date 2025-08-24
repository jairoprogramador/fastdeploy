package port

import "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entity"

type Repository interface {
	Exists() (bool, error)
	Load() (entity.Configuration, error)
	Save(config entity.Configuration) error
}
