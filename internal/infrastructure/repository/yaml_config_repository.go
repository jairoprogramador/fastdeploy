package repository

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/domain/config/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/config/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/yaml"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

const (
	erroConfigNotFound = "file config not found in %s"

	msgSuccessFileConfigExists = "file config exists in %s"
	msgSuccessSaveConfig       = "save file %s successful"
)

type yamlConfigRepository struct {
	yamlRepository yaml.YamlController
	fileRepository file.FileController
	router         port.PathService
	fileLogger     *logger.FileLogger
}

func NewYamlConfigRepository(
	yamlRepository yaml.YamlController,
	fileRepository file.FileController,
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

func (r *yamlConfigRepository) Load() result.InfraResult {
	path := r.router.GetFullPathConfigFile()

	if result := r.exists(path); !result.IsSuccess() {
		return result
	}

	var configEntity *entity.ConfigEntity
	if err := r.yamlRepository.Load(path, &configEntity); err != nil {
		return result.NewError(err)
	}

	return result.NewResult(&configEntity)
}

func (r *yamlConfigRepository) Save(config *entity.ConfigEntity) result.InfraResult {
	path := r.router.GetFullPathConfigFile()

	if response := r.exists(path); response.IsSuccess() {
		if err := r.fileRepository.DeleteFile(path); err != nil {
			return result.NewError(err)
		}
	}

	if err := r.yamlRepository.Save(path, config); err != nil {
		return result.NewError(err)
	}

	return result.NewResult(fmt.Sprintf(msgSuccessSaveConfig, path))
}

func (r *yamlConfigRepository) exists(path string) result.InfraResult {
	exists, err := r.fileRepository.ExistsFile(path)
	if err != nil {
		return result.NewError(err)
	}
	if !exists {
		return r.logError(fmt.Errorf(erroConfigNotFound, path))
	}
	return result.NewResult(fmt.Sprintf(msgSuccessFileConfigExists, path))
}

func (r *yamlConfigRepository) logError(err error) result.InfraResult {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return result.NewError(err)
}
