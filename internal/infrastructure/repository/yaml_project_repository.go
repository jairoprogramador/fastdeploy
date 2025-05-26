package repository

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model/logger"
	"github.com/jairoprogramador/fastdeploy/internal/domain/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	erroProjectNotFound  = "file project not found in %s"
	erroResourceNotFound = "the resource was not found in the directory %s"

	msgSuccessFileProjectExists = "file project exists in %s"
	msgSuccessSaveProject       = "save file %s successful"
)

type yamlProjectRepository struct {
	yamlRepository adapter.YamlController
	fileRepository adapter.FileController
	router         port.PathService
	fileLogger     *logger.FileLogger
}

func NewYamlProjectRepository(
	yamlRepository adapter.YamlController,
	fileRepository adapter.FileController,
	router port.PathService,
	fileLogger *logger.FileLogger,
) repository.ProjectRepository {
	return &yamlProjectRepository{
		yamlRepository: yamlRepository,
		fileRepository: fileRepository,
		router:         router,
		fileLogger:     fileLogger,
	}
}

func (r *yamlProjectRepository) Load() model.InfraResultEntity {
	path := r.router.GetPathProjectFile()

	if result := r.exists(path); !result.IsSuccess() {
		return result
	}

	var project model.ProjectEntity
	if err := r.yamlRepository.Load(path, &project); err != nil {
		return model.NewError(err)
	}

	return model.NewResult(&project)
}

func (r *yamlProjectRepository) Save(project *model.ProjectEntity) model.InfraResultEntity {
	path := r.router.GetPathProjectFile()

	if result := r.exists(path); result.IsSuccess() {
		if err := r.fileRepository.DeleteFile(path); err != nil {
			return model.NewError(err)
		}
	}

	if err := r.yamlRepository.Save(path, project); err != nil {
		return model.NewError(err)
	}

	return model.NewResult(fmt.Sprintf(msgSuccessSaveProject, path))
}

func (r *yamlProjectRepository) GetName() model.InfraResultEntity {
	pathWorkingDir, err := os.Getwd()
	if err != nil {
		return r.logError(err)
	}

	return model.NewResult(filepath.Base(pathWorkingDir))
}

func (r *yamlProjectRepository) GetFullPathResource() model.InfraResultEntity {
	directory := "target"
	existsDir, err := r.fileRepository.ExistsDirectory(directory)
	if err == nil && existsDir {
		fullPathJarFiles, err := r.getFullPathResources(directory)
		if err == nil {
			return model.NewResult(fullPathJarFiles[0])
		}
		return r.logError(err)
	}
	return r.logError(err)
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
		return nil, fmt.Errorf(erroResourceNotFound, pathDirectory)
	}

	return pathFiles, nil
}

func (r *yamlProjectRepository) exists(path string) model.InfraResultEntity {
	exists, err := r.fileRepository.ExistsFile(path)
	if err != nil {
		return model.NewError(err)
	}
	if !exists {
		return r.logError(fmt.Errorf(erroProjectNotFound, path))
	}
	return model.NewResult(fmt.Sprintf(msgSuccessFileProjectExists, path))
}

func (r *yamlProjectRepository) logError(err error) model.InfraResultEntity {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return model.NewError(err)
}
