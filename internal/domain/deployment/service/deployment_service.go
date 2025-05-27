package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
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
		s.setDefaultType(deployment)
		s.setCheckContainer(deployment)
	}

	return &model.DeploymentEntity{}, result.Error
}

func (s *deploymentService) setDefaultType(deployment *model.DeploymentEntity) {
	for i := range deployment.Steps {
		if deployment.Steps[i].Type == "" {
			deployment.Steps[i].Type = string(model.Command)
		}
	}
}

func (s *deploymentService) setCheckContainer(deployment *model.DeploymentEntity) {
	if deployment.HasType(string(model.Container)) {
		setupStep := model.Step{
			Name:    string(model.Check),
			Type:    string(model.Check),
			Timeout: "30s",
			Then:    validator.ThenFinish,
		}
		deployment.Steps = append([]model.Step{setupStep}, deployment.Steps...)
	}
}
