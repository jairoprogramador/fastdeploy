package repository

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/template"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/yaml"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
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

	if err := r.ensureDeploymentFileExists(path); err != nil {
		return result.NewError(err)
	}

	var deployment model.DeploymentEntity

	if err := r.yamlPort.Load(path, &deployment); err != nil {
		return result.NewError(err)
	}

	return result.NewResult(deployment)
}

func (r *yamlDeploymentRepository) ensureDeploymentFileExists(pathTemplate string) error {
	exists, err := r.filePort.ExistsFile(pathTemplate)
	if err != nil {
		return err
	}

	if !exists {
		if err := r.filePort.WriteFile(pathTemplate, template.DeploymentTemplate); err != nil {
			return err
		}
	}

	return nil
}

func (r *yamlDeploymentRepository) logError(err error) result.InfraResult {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return result.NewError(err)
}
