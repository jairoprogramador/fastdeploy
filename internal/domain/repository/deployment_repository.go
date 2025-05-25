package repository

import (
	"deploy/internal/domain/engine/model"
)

type DeploymentRepository interface {
	Load() (*model.DeploymentEntity, error)
}
