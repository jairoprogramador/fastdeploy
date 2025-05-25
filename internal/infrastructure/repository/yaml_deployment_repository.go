package repository

import (
	"deploy/internal/domain/engine/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/service"
	"deploy/internal/infrastructure/adapter"
)

type yamlDeploymentRepository struct {
	yamlRepository adapter.YamlController
	fileRepository adapter.FileController
	router         *service.PathService
}

func NewYamlDeploymentRepository(
	yamlRepository adapter.YamlController,
	fileRepository adapter.FileController,
	router *service.PathService,
) repository.DeploymentRepository {
	return &yamlDeploymentRepository{
		yamlRepository: yamlRepository,
		fileRepository: fileRepository,
		router:         router,
	}
}

func (r *yamlDeploymentRepository) Load() (*model.DeploymentEntity, error) {
	path := r.router.GetFullPathDeploymentFile()

	exists, err := r.fileRepository.ExistsFile(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, service.ErrDeploymentNotFound
	}

	var deployment *model.DeploymentEntity
	response := r.yamlRepository.Load(path, &deployment)
	if !response.IsSuccess() {
		return &model.DeploymentEntity{}, response.Error
	}

	return deployment, nil
}
