package repository

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model/logger"
	"github.com/jairoprogramador/fastdeploy/internal/domain/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter"
	"fmt"
)

const (
	erroConfigNotFound = "file config not found in %s"

	msgSuccessFileConfigExists = "file config exists in %s"
	msgSuccessSaveConfig       = "save file %s successful"
)

type yamlConfigRepository struct {
	yamlRepository adapter.YamlController
	fileRepository adapter.FileController
	router         port.PathService
	fileLogger     *logger.FileLogger
}

func NewYamlConfigRepository(
	yamlRepository adapter.YamlController,
	fileRepository adapter.FileController,
	router port.PathService,
	fileLogger *logger.FileLogger,
) repository.ConfigRepository {
	return &yamlConfigRepository{
		yamlRepository: yamlRepository,
		fileRepository: fileRepository,
		router:         router,
		fileLogger:     fileLogger,
	}
}

func (r *yamlConfigRepository) Load() model.InfraResultEntity {
	path := r.router.GetFullPathConfigFile()

	if result := r.exists(path); !result.IsSuccess() {
		return result
	}

	var configEntity *model.ConfigEntity
	if err := r.yamlRepository.Load(path, &configEntity); err != nil {
		return model.NewError(err)
	}

	return model.NewResult(&configEntity)
}

func (r *yamlConfigRepository) Save(config *model.ConfigEntity) model.InfraResultEntity {
	path := r.router.GetFullPathConfigFile()

	if result := r.exists(path); result.IsSuccess() {
		if err := r.fileRepository.DeleteFile(path); err != nil {
			return model.NewError(err)
		}
	}

	if err := r.yamlRepository.Save(path, config); err != nil {
		return model.NewError(err)
	}

	return model.NewResult(fmt.Sprintf(msgSuccessSaveConfig, path))
}

func (r *yamlConfigRepository) exists(path string) model.InfraResultEntity {
	exists, err := r.fileRepository.ExistsFile(path)
	if err != nil {
		return model.NewError(err)
	}
	if !exists {
		return r.logError(fmt.Errorf(erroConfigNotFound, path))
	}
	return model.NewResult(fmt.Sprintf(msgSuccessFileConfigExists, path))
}

func (r *yamlConfigRepository) logError(err error) model.InfraResultEntity {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return model.NewError(err)
}
