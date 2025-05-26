package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
	"github.com/jairoprogramador/fastdeploy/internal/domain/repository"
)

type DeploymentService interface {
	Load() (*model.DeploymentEntity, error)
}

type deploymentService struct {
	deploymentRepository repository.DeploymentRepository
}

func NewDeploymentService(
	deploymentRepository repository.DeploymentRepository,
) DeploymentService {
	return &deploymentService{
		deploymentRepository: deploymentRepository,
	}
}

func (s *deploymentService) Load() (*model.DeploymentEntity, error) {
	result := s.deploymentRepository.Load()
	if result.IsSuccess() {
		deployment := result.Result.(*model.DeploymentEntity)
		for i := range deployment.Steps {
			if deployment.Steps[i].Type == "" {
				deployment.Steps[i].Type = string(model.Command)
			}
		}

		if deployment.HasType(string(model.Container)) {
			setupStep := model.Step{
				Name:    string(model.Setup),
				Type:    string(model.Setup),
				Timeout: "30s",
				Then:    validator.ThenFinish,
			}
			deployment.Steps = append([]model.Step{setupStep}, deployment.Steps...)
		}
	}

	return &model.DeploymentEntity{}, result.Error
}
