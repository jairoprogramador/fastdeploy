package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/config"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/git"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
)

type ProjectFactory struct{}

func NewProjectFactory() *ProjectFactory {
	return &ProjectFactory{}
}

func (pf *ProjectFactory) CreateRepository() domain.ProjectRepository {
	fileSystem := &filesystem.OSFileSystem{}

	pathResolver := pf.CreatePathResolver()
	serializer := NewYAMLProjectSerializer()

	return NewYAMLProjectRepository(fileSystem, pathResolver, serializer)
}

func (pf *ProjectFactory) CreateService() domain.ProjectService {
	repository := pf.CreateRepository()
	validator := domain.NewProjectValidator()

	return domain.NewProjectService(repository, validator)
}

func (pf *ProjectFactory) CreatePathResolver() ProjectPathResolver {
	workingDir := &filesystem.OSWorkingDirectory{}
	return NewProjectPathResolver(workingDir)
}

func (pf *ProjectFactory) CreateCreator() domain.ProjectCreator {
	pathResolver := pf.CreatePathResolver()
	idGenerator := NewProjectIDGenerator()

	return NewProjectCreator(pathResolver, idGenerator)
}

// CreateDefaultProjectInitializerService crea un inicializador de proyecto con dependencias por defecto
func (pf *ProjectFactory) CreateInitialize() domain.ProjectInitialize {
	projectService := pf.CreateService()
	projectCreate := pf.CreateCreator()
	configService := config.NewConfigFactory().CreateService()
	gitService := git.NewGitFactory().CreateService()
	initializerService := domain.NewProjectInitialize(projectService, projectCreate, configService, gitService)

	return initializerService
}
