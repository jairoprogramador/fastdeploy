package repository

import (
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/service"
	"deploy/internal/domain/service/router"
	"deploy/internal/infrastructure/adapter"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type yamlProjectRepository struct {
	yamlRepository adapter.YamlRepository
	fileRepository adapter.FileRepository
	router         *router.Router
}

// NewYamlProjectRepository creates a new instance of ProjectRepository
func NewYamlProjectRepository(
	yamlRepository adapter.YamlRepository,
	fileRepository adapter.FileRepository,
	router *router.Router,
) repository.ProjectRepository {
	return &yamlProjectRepository{
		yamlRepository: yamlRepository,
		fileRepository: fileRepository,
		router:         router,
	}
}

// Load loads a project from the file system
func (r *yamlProjectRepository) Load() (*model.ProjectEntity, error) {
	path := r.router.GetPathProjectFile()

	if err := r.exists(path); err != nil {
		return &model.ProjectEntity{}, err
	}

	var project model.ProjectEntity
	response := r.yamlRepository.Load(path, &project)
	if !response.IsSuccess() {
		return &model.ProjectEntity{}, response.Error
	}
	return &project, nil
}

// Save saves a project to the file system
func (r *yamlProjectRepository) Save(project *model.ProjectEntity) error {
	path := r.router.GetPathProjectFile()

	if err := r.exists(path); err == nil {
		if err := r.fileRepository.DeleteFile(path); err != nil {
			return err
		}
	}

	if response := r.yamlRepository.Save(path, project); !response.IsSuccess() {
		return response.Error
	}

	return nil
}

func (r *yamlProjectRepository) GetName() (string, error) {
	pathWorkingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(pathWorkingDir), nil
}

func (r *yamlProjectRepository) GetFullPathResource() (string, error) {
	directory := "target"
	existsDir, err := r.fileRepository.ExistsDirectory(directory)
	if err == nil && existsDir {
		fullPathJarFiles, err := r.getFullPathResources(directory)
		if err == nil {
			return fullPathJarFiles[0], nil
		}
		return "", err
	}
	return "", err
}

func (r *yamlProjectRepository) getFullPathResources(pathDirectory string) ([]string, error) {
	files, err := r.fileRepository.ReadDirectory(pathDirectory)
	if err != nil {
		return nil, err
	}

	var pathFiles []string

	for _, archivo := range files {
		if !archivo.IsDir() && strings.HasSuffix(archivo.Name(), ".jar") &&
			!strings.Contains(archivo.Name(), "sources") &&
			!strings.Contains(archivo.Name(), "original") {

			pathFileRelative := r.fileRepository.GetPath(pathDirectory, archivo.Name())
			pathFileAbsolute, err := r.fileRepository.GetPathAbsolute(pathFileRelative)
			if err != nil {
				return nil, err
			}
			pathFiles = append(pathFiles, pathFileAbsolute)
		}
	}

	if len(pathFiles) <= 0 {
		return nil, fmt.Errorf("the resource was not found in the directory %s", pathDirectory)
	}

	return pathFiles, nil
}

func (r *yamlProjectRepository) exists(path string) error {
	exists, err := r.fileRepository.ExistsFile(path)
	if err != nil {
		return err
	}
	if !exists {
		return service.ErrProjectNotFound
	}
	return nil
}
