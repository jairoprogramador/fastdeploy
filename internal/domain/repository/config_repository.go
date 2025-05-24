package repository

import (
	"deploy/internal/domain/model"
)

type ConfigRepository interface {
	Load() (*model.ConfigEntity, error)
	Save(config *model.ConfigEntity) error
}
