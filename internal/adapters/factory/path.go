package factory

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/config"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/git"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/project"
)

type PathFactory interface {
	CreateProjectPathResolver() project.ProjectPathResolver
	CreateConfigPathResolver() config.ConfigPathResolver
	CreateGitPathResolver() git.GitPathResolver
}
