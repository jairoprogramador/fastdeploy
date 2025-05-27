package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
)

type DeploymentService interface {
	Load() (*entity.DeploymentEntity, error)
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

func (s *deploymentService) Load() (*entity.DeploymentEntity, error) {
	result := s.deploymentRepository.Load()
	if result.IsSuccess() {
		deployment := result.Result.(*entity.DeploymentEntity)
		for i := range deployment.Steps {
			if deployment.Steps[i].Type == "" {
				deployment.Steps[i].Type = string(entity.Command)
			}
		}

		if deployment.HasType(string(entity.Container)) {
			setupStep := entity.Step{
				Name:    string(entity.Setup),
				Type:    string(entity.Setup),
				Timeout: "30s",
				Then:    validator.ThenFinish,
			}
			deployment.Steps = append([]entity.Step{setupStep}, deployment.Steps...)
		}
	}

	return &entity.DeploymentEntity{}, result.Error
}
