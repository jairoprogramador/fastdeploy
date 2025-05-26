package repository

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
)

type ProjectRepository interface {
	Load() model.InfraResultEntity
	Save(project *model.ProjectEntity) model.InfraResultEntity
	GetName() model.InfraResultEntity
	GetFullPathResource() model.InfraResultEntity
}
