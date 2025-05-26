package repository

import (
	modelEngine "github.com/jairoprogramador/fastdeploy/internal/domain/engine/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model/logger"
	"github.com/jairoprogramador/fastdeploy/internal/domain/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter"
	"fmt"
)

const (
	erroDeploymentNotFound = "file deployment not found in %s"

	msgSuccessFileDeploymentExists = "file deployment exists in %s"
)

type yamlDeploymentRepository struct {
	yamlRepository adapter.YamlController
	fileRepository adapter.FileController
	router         port.PathService
	fileLogger     *logger.FileLogger
}

func NewYamlDeploymentRepository(
	yamlRepository adapter.YamlController,
	fileRepository adapter.FileController,
	router port.PathService,
	fileLogger *logger.FileLogger,
) repository.DeploymentRepository {
	return &yamlDeploymentRepository{
		yamlRepository: yamlRepository,
		fileRepository: fileRepository,
		router:         router,
		fileLogger:     fileLogger,
	}
}

func (r *yamlDeploymentRepository) Load() model.InfraResultEntity {
	path := r.router.GetFullPathDeploymentFile()

	if result := r.exists(path); !result.IsSuccess() {
		return result
	}

	var deployment *modelEngine.DeploymentEntity
	if err := r.yamlRepository.Load(path, &deployment); err != nil {
		return model.NewError(err)
	}

	return model.NewResult(&deployment)
}

func (r *yamlDeploymentRepository) exists(path string) model.InfraResultEntity {
	exists, err := r.fileRepository.ExistsFile(path)
	if err != nil {
		return model.NewError(err)
	}
	if !exists {
		return r.logError(fmt.Errorf(erroDeploymentNotFound, path))
	}
	return model.NewResult(fmt.Sprintf(msgSuccessFileDeploymentExists, path))
}

func (r *yamlDeploymentRepository) logError(err error) model.InfraResultEntity {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return model.NewError(err)
}
