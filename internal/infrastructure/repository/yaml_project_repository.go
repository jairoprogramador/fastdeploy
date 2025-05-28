package repository

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/repository"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/yaml"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
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
	yamlPort   yaml.YamlPort
	filePort   file.FilePort
	pathPort   port.PathPort
	fileLogger *logger.FileLogger
}

func NewProjectRepository(
	yamlPort yaml.YamlPort,
	filePort file.FilePort,
	pathPort port.PathPort,
	fileLogger *logger.FileLogger,
) repository.ProjectRepository {
	return &yamlProjectRepository{
		yamlPort:   yamlPort,
		filePort:   filePort,
		pathPort:   pathPort,
		fileLogger: fileLogger,
	}
}

func (r *yamlProjectRepository) Load() result.InfraResult {
	path := r.pathPort.GetPathProjectFile()

	if _, err := r.exists(path); err != nil {
		return r.logError(err)
	}

	var project model.ProjectEntity
	if err := r.yamlPort.Load(path, &project); err != nil {
		return result.NewError(err)
	}

	return result.NewResult(project)
}

func (r *yamlProjectRepository) Save(project *model.ProjectEntity) result.InfraResult {
	path := r.pathPort.GetPathProjectFile()

	if exists, _ := r.exists(path); exists {
		if err := r.filePort.DeleteFile(path); err != nil {
			return result.NewError(err)
		}
	}

	if err := r.yamlPort.Save(path, project); err != nil {
		return result.NewError(err)
	}

	return result.NewResult(fmt.Sprintf(msgSuccessSaveProject, path))
}

func (r *yamlProjectRepository) GetName() result.InfraResult {
	pathWorkingDir, err := os.Getwd()
	if err != nil {
		return r.logError(err)
	}

	return result.NewResult(filepath.Base(pathWorkingDir))
}

func (r *yamlProjectRepository) GetFullPathResource() result.InfraResult {
	directory := "target"
	existsDir, err := r.filePort.ExistsDirectory(directory)
	if err == nil && existsDir {
		fullPathJarFiles, err := r.getFullPathResources(directory)
		if err == nil {
			return result.NewResult(fullPathJarFiles[0])
		}
		return r.logError(err)
	}
	return r.logError(err)
}

func (r *yamlProjectRepository) getFullPathResources(pathDirectory string) ([]string, error) {
	files, err := r.filePort.ReadDirectory(pathDirectory)
	if err != nil {
		return nil, err
	}

	var pathFiles []string

	for _, archivo := range files {
		if !archivo.IsDir() && strings.HasSuffix(archivo.Name(), ".jar") &&
			!strings.Contains(archivo.Name(), "sources") &&
			!strings.Contains(archivo.Name(), "original") {

			pathFileRelative := r.filePort.GetPath(pathDirectory, archivo.Name())
			pathFileAbsolute, err := r.filePort.GetPathAbsolute(pathFileRelative)
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

func (r *yamlProjectRepository) exists(path string) (bool, error) {
	exists, err := r.filePort.ExistsFile(path)
	if !exists {
		return exists, fmt.Errorf(erroProjectNotFound, path)
	}
	return exists, err
}

func (r *yamlProjectRepository) logError(err error) result.InfraResult {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return result.NewError(err)
}
