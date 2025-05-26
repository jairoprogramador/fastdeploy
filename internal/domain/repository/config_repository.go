package repository

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
)

type ConfigRepository interface {
	Load() model.InfraResultEntity
	Save(config *model.ConfigEntity) model.InfraResultEntity
}
