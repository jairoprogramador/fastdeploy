package repository

import "deploy/internal/domain/model"

type ProjectRepository interface {
	Load(path string) (model.Project, error)
	Create(path string, project *model.Project) error
	GetProjectName() (string, error)
} 
