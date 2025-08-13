package impl

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/config"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/factory"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/git"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/project"
)

type PathFactoryImpl struct{}

func NewPathFactory() factory.PathFactory {
	return &PathFactoryImpl{}
}

func (pf *PathFactoryImpl) CreateProjectPathResolver() project.ProjectPathResolver {
	workingDir := &filesystem.OSWorkingDirectory{}
	return project.NewProjectPathResolver(workingDir)
}

func (pf *PathFactoryImpl) CreateConfigPathResolver() config.ConfigPathResolver {
	userSystem := &filesystem.OSUserSystem{}
	return config.NewConfigPathResolver(userSystem)
}

func (pf *PathFactoryImpl) CreateGitPathResolver() git.GitPathResolver {
	configPathResolver := pf.CreateConfigPathResolver()
	return git.NewGitPathResolver(configPathResolver)
}
