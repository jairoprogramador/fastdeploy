package repository

import (
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/service"
	"deploy/internal/infrastructure/adapter"
)

type yamlConfigRepository struct {
	yamlRepository adapter.YamlController
	fileRepository adapter.FileController
	router         *service.PathService
}

// NewYamlConfigRepository creates a new instance of ConfigRepository
func NewYamlConfigRepository(
	yamlRepository adapter.YamlController,
	fileRepository adapter.FileController,
	router *service.PathService,
) repository.ConfigRepository {
	return &yamlConfigRepository{
		yamlRepository: yamlRepository,
		fileRepository: fileRepository,
		router:         router,
	}
}

// Load loads the configuration from storage
func (r *yamlConfigRepository) Load() (*model.ConfigEntity, error) {
	path := r.router.GetFullPathGlobalConfigFile()

	if err := r.exists(path); err != nil {
		return &model.ConfigEntity{}, err
	}

	var configEntity model.ConfigEntity
	response := r.yamlRepository.Load(path, &configEntity)
	if !response.IsSuccess() {
		return &model.ConfigEntity{}, response.Error
	}

	return &configEntity, nil
}

// Save saves the configuration to storage
func (r *yamlConfigRepository) Save(globalConfig *model.ConfigEntity) error {
	path := r.router.GetFullPathGlobalConfigFile()

	if err := r.exists(path); err == nil {
		if err := r.fileRepository.DeleteFile(path); err != nil {
			return err
		}
	}

	if response := r.yamlRepository.Save(path, globalConfig); !response.IsSuccess() {
		return response.Error
	}

	return nil
}

func (r *yamlConfigRepository) exists(path string) error {
	exists, err := r.fileRepository.ExistsFile(path)
	if err != nil {
		return err
	}
	if !exists {
		return service.ErrConfigNotFound
	}
	return nil
}
