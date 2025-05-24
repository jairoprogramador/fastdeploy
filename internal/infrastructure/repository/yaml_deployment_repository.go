package repository

import (
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/service"
	"deploy/internal/domain/service/router"
	"deploy/internal/infrastructure/adapter"
)

type yamlDeploymentRepository struct {
	yamlRepository adapter.YamlRepository
	fileRepository adapter.FileRepository
	router         *router.Router
}

func NewYamlDeploymentRepository(
	yamlRepository adapter.YamlRepository,
	fileRepository adapter.FileRepository,
	router *router.Router,
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
