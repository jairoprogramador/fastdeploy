package repository

import (
	"fmt"
	modelEngine "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
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
	yamlPort   yaml.YamlPort
	filePort   file.FilePort
	pathPort   port.PathPort
	fileLogger *logger.FileLogger
}

func NewDeploymentRepository(
	yamlPort yaml.YamlPort,
	filePort file.FilePort,
	pathPort port.PathPort,
	fileLogger *logger.FileLogger,
) repository.DeploymentRepository {
	return &yamlDeploymentRepository{
		yamlPort:   yamlPort,
		filePort:   filePort,
		pathPort:   pathPort,
		fileLogger: fileLogger,
	}
}

func (r *yamlDeploymentRepository) Load() result.InfraResult {
	path := r.pathPort.GetFullPathDeploymentFile()

	if result := r.exists(path); !result.IsSuccess() {
		return result
	}

	var deployment *modelEngine.DeploymentEntity
	if err := r.yamlPort.Load(path, &deployment); err != nil {
		return result.NewError(err)
	}

	return result.NewResult(&deployment)
}

func (r *yamlDeploymentRepository) exists(path string) result.InfraResult {
	exists, err := r.filePort.ExistsFile(path)
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
