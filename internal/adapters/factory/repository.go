package factory

import (
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
)

type RepositoryFactory interface {
	CreateProjectRepository() project.ProjectRepository
	CreateConfigRepository() config.ConfigRepository
}
