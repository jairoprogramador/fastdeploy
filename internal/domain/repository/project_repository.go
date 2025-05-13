package repository

import "deploy/internal/domain/model"

type ProjectRepository interface {
	Load() (model.Project, error)
	Create(project *model.Project) error
	RemoveFile() error
	GetProjectName() (string, error)
} 
