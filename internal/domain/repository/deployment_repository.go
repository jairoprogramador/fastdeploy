package repository

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
)

type DeploymentRepository interface {
	Load() model.InfraResultEntity
}
