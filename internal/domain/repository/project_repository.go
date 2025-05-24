package repository

import (
	"deploy/internal/domain/model"
)

type ProjectRepository interface {
	Load() (*model.ProjectEntity, error)
	Save(project *model.ProjectEntity) error
	GetName() (string, error)
	GetFullPathResource() (string, error)
}
