package repository

import (
	"sync"
	"deploy/internal/infrastructure/filesystem"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/template"
	"deploy/internal/domain/model"
	"deploy/internal/domain"
)

type projectRepositoryImpl struct {}

var (
    instanceProjectRepository     repository.ProjectRepository
    instanceOnceProjectRepository sync.Once
)

func GetProjectRepository() repository.ProjectRepository {
    instanceOnceProjectRepository.Do(func() {
        instanceProjectRepository = &projectRepositoryImpl{}
    })
    return instanceProjectRepository
}

func (s *projectRepositoryImpl) IsInitialized(rootDirectory, nameProjectFile string) bool {
	filePath := filesystem.GetPath(rootDirectory, nameProjectFile)
	exists := filesystem.FileExists(filePath)
	if !exists {
		return exists
	}

	filePath = constants.DockerfileTemplateFilePath
	exists = filesystem.FileExists(filePath)
	if !exists {
		return exists
	}

	filePath = constants.DockercomposeTemplateFilePath
    return filesystem.FileExists(filePath)
}

func (s *projectRepositoryImpl) GetOrganizationName() string {
    return ""
}

func (s *projectRepositoryImpl) GetProjectId() (string, error) {
	return filesystem.GetParentDirectory()
}

func (s *projectRepositoryImpl) GetTeamName() string {
    return ""
}

func (st *projectRepositoryImpl) Load() (model.Project, error) {
	filePath := filesystem.GetPath(constants.RootDirectory, constants.NameProjectFile)
	return filesystem.LoadFromYAML[model.Project](filePath)
}

func (s *projectRepositoryImpl) SaveDockercomposeTemplate() *model.Response {
	filePath := constants.DockercomposeTemplateFilePath
	if err := filesystem.CreateDirectoryFilePath(filePath); err != nil {
		return model.GetNewResponseError(err)
	}
	err := filesystem.WriteFile(filePath, template.ComposeTemplate)
	if err != nil {
		return model.GetNewResponseError(err)
	}
	return model.GetNewResponse()
}

func (s *projectRepositoryImpl) SaveDockerfileTemplate() *model.Response {
	filePath := constants.DockerfileTemplateFilePath
	if err := filesystem.CreateDirectoryFilePath(filePath); err != nil {
		return model.GetNewResponseError(err)
	}
	err := filesystem.WriteFile(filePath, template.DockerfileTemplate)
	if err != nil {
		return model.GetNewResponseError(err)
	}
	return model.GetNewResponse()
}

func (st projectRepositoryImpl) Save(project *model.Project) *model.Response {
	filePath := filesystem.GetPath(constants.RootDirectory, constants.NameProjectFile)
	err := filesystem.SaveToYAML(project, filePath)
	if err != nil {
		return model.GetNewResponseError(err)
	}
	return model.GetNewResponse()
}
