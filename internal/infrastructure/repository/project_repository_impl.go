// Package repository implementa el patrón Repository para el manejo de proyectos
package repository

import (
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/filesystem"
	"sync"

	//"deploy/internal/domain/template"
	"deploy/internal/domain/model"
	//"deploy/internal/domain"
)

// projectRepositoryImpl implementa la interfaz ProjectRepository
type projectRepositoryImpl struct{}

var (
	instanceProjectRepository     repository.ProjectRepository
	instanceOnceProjectRepository sync.Once
)

// NewProjectRepository crea una nueva instancia del repositorio de proyectos
// utilizando el patrón Singleton
func GetProjectRepository() repository.ProjectRepository {
	instanceOnceProjectRepository.Do(func() {
		instanceProjectRepository = &projectRepositoryImpl{}
	})
	return instanceProjectRepository
}

// Load carga el proyecto desde el archivo YAML
func (st *projectRepositoryImpl) Load() (model.Project, error) {
	filePath := st.getPathProjectFile()
	exists, err := filesystem.ExistsFile(filePath)
	if !exists {
		return model.Project{}, err
	}
	return filesystem.LoadFromYAML[model.Project](filePath)
}

// RemoveFile elimina el archivo del proyecto
func (st *projectRepositoryImpl) RemoveFile() error {
	filePath := st.getPathProjectFile()
	return filesystem.RemoveFile(filePath)
}

// Create crea un nuevo proyecto y lo guarda en el archivo YAML
func (st *projectRepositoryImpl) Create(project *model.Project) error {
	if err := st.RemoveFile(); err != nil {
		return err
	}

	filePath := st.getPathProjectFile()
	err := filesystem.SaveToYAML(project, filePath)
	if err != nil {
		return err
	}
	return nil
}

// GetProjectName obtiene el nombre del proyecto desde el directorio actual
func (s *projectRepositoryImpl) GetProjectName() (string, error) {
	return filesystem.GetParentDirectory()
}

// getPathProjectFile obtiene la ruta del archivo del proyecto
func (s *projectRepositoryImpl) getPathProjectFile() string {
	return filesystem.GetPath(ProjectDirectory, ProjectFile)
}

/* func (s *projectRepositoryImpl) Exists() bool {
	filePath := filesystem.GetPath(constants.RootDirectory, constants.NameProjectFile)
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
} */

/* func (s *projectRepositoryImpl) SaveDockercomposeTemplate() *model.Response {
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
} */

/* func (st projectRepositoryImpl) Save(project *model.Project) *model.Response {
	filePath := filesystem.GetPath(constants.RootDirectory, constants.NameProjectFile)
	err := filesystem.SaveToYAML(project, filePath)
	if err != nil {
		return model.GetNewResponseError(err)
	}
	return model.GetNewResponse()
} */
