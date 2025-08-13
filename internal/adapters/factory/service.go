package factory

import (
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/git"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
)

type ServiceFactory interface {
	CreateProjectService() project.ProjectService
	CreateConfigService() config.ConfigService
	CreateGitService() git.GitService
}
