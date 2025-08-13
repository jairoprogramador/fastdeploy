package impl

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/config"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/factory"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/project"
	configDomain "github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	projectDomain "github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
)

type RepositoryFactoryImpl struct{}

func NewRepositoryFactory() factory.RepositoryFactory {
	return &RepositoryFactoryImpl{}
}

func (rf *RepositoryFactoryImpl) CreateProjectRepository() projectDomain.ProjectRepository {
	fileSystem := &filesystem.OSFileSystem{}
	workingDir := &filesystem.OSWorkingDirectory{}

	pathResolver := project.NewProjectPathResolver(workingDir)
	serializer := project.NewYAMLProjectSerializer()

	return project.NewYAMLProjectRepository(fileSystem, pathResolver, serializer)
}

func (rf *RepositoryFactoryImpl) CreateConfigRepository() configDomain.ConfigRepository {
	userSystem := &filesystem.OSUserSystem{}
	fileSystem := &filesystem.OSFileSystem{}

	pathResolver := config.NewConfigPathResolver(userSystem)
	serializer := config.NewYAMLConfigSerializer()

	return config.NewYAMLConfigRepository(fileSystem, userSystem, pathResolver, serializer)
}
