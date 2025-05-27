package repository

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

type ProjectRepository interface {
	Load() result.InfraResult
	Save(project *entity.ProjectEntity) result.InfraResult
	GetName() result.InfraResult
	GetFullPathResource() result.InfraResult
}
