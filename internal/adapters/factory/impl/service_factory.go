package impl

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/config"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/factory"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/git"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/project"
	configDomain "github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	gitDomain "github.com/jairoprogramador/fastdeploy/internal/core/domain/git"
	projectDomain "github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
)

type ServiceFactoryImpl struct{}

func NewServiceFactory() factory.ServiceFactory {
	return &ServiceFactoryImpl{}
}

func (sf *ServiceFactoryImpl) CreateProjectService() projectDomain.ProjectService {
	fileSystem := &filesystem.OSFileSystem{}
	workingDir := &filesystem.OSWorkingDirectory{}

	pathResolver := project.NewProjectPathResolver(workingDir)
	serializer := project.NewYAMLProjectSerializer()

	repository := project.NewYAMLProjectRepository(fileSystem, pathResolver, serializer)
	validator := projectDomain.NewProjectValidator()

	return projectDomain.NewProjectService(repository, validator)
}

func (sf *ServiceFactoryImpl) CreateConfigService() configDomain.ConfigService {
	userSystem := &filesystem.OSUserSystem{}
	fileSystem := &filesystem.OSFileSystem{}

	pathResolver := config.NewConfigPathResolver(userSystem)
	serializer := config.NewYAMLConfigSerializer()
	repository := config.NewYAMLConfigRepository(fileSystem, userSystem, pathResolver, serializer)
	validator := configDomain.NewConfigValidatorImpl()

	return configDomain.NewConfigService(repository, validator)
}

func (sf *ServiceFactoryImpl) CreateGitService() gitDomain.GitService {
	pathResolver := NewPathFactory().CreateGitPathResolver()
	return git.NewGitService(pathResolver)
}
