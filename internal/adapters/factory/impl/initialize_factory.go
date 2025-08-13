package impl

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/factory"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/project"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
)

type InitializeFactoryImpl struct{}

func NewInitializeFactory() factory.InitializeFactory {
	return &InitializeFactoryImpl{}
}

func (pf *InitializeFactoryImpl) CreateInitialize() domain.ProjectInitialize {
	pathResolver := NewPathFactory().CreateProjectPathResolver()
	idGenerator := project.NewProjectIDGenerator()
	projectCreate := project.NewProjectCreator(pathResolver, idGenerator)

	serviceFactory := NewServiceFactory()
	projectService := serviceFactory.CreateProjectService()
	configService := serviceFactory.CreateConfigService()

	gitService := serviceFactory.CreateGitService()
	initializerService := domain.NewProjectInitialize(projectService, projectCreate, configService, gitService)

	return initializerService
}
