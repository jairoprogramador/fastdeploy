package repository

import (
	"fmt"
	modelEngine "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/yaml"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

const (
	erroDeploymentNotFound = "file deployment not found in %s"

	msgSuccessFileDeploymentExists = "file deployment exists in %s"
)

type yamlDeploymentRepository struct {
	yamlRepository yaml.YamlController
	fileRepository file.FileController
	router         port.PathService
	fileLogger     *logger.FileLogger
}

func NewYamlDeploymentRepository(
	yamlRepository yaml.YamlController,
	fileRepository file.FileController,
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

func (r *yamlDeploymentRepository) Load() result.InfraResult {
	path := r.router.GetFullPathDeploymentFile()

	if result := r.exists(path); !result.IsSuccess() {
		return result
	}

	var deployment *modelEngine.DeploymentEntity
	if err := r.yamlRepository.Load(path, &deployment); err != nil {
		return result.NewError(err)
	}

	return result.NewResult(&deployment)
}

func (r *yamlDeploymentRepository) exists(path string) result.InfraResult {
	exists, err := r.fileRepository.ExistsFile(path)
	if err != nil {
		return result.NewError(err)
	}
	if !exists {
		return r.logError(fmt.Errorf(erroDeploymentNotFound, path))
	}
	return result.NewResult(fmt.Sprintf(msgSuccessFileDeploymentExists, path))
}

func (r *yamlDeploymentRepository) logError(err error) result.InfraResult {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return result.NewError(err)
}
