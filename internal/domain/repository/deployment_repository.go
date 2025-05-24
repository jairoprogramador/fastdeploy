package repository

import (
	"deploy/internal/domain/model"
)

type DeploymentRepository interface {
	Load() (*model.DeploymentEntity, error)
}
