package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/service"
)

type DeploymentService interface {
	Load() (*model.DeploymentEntity, error)
}

type deploymentService struct {
	deploymentRepository repository.DeploymentRepository
	storeService         service.StoreServicePort
}

func NewDeploymentService(
	deploymentRepository repository.DeploymentRepository,
	storeService service.StoreServicePort,
) DeploymentService {
	return &deploymentService{
		deploymentRepository: deploymentRepository,
		storeService:         storeService,
	}
}

func (s *deploymentService) Load() (*model.DeploymentEntity, error) {
	result := s.deploymentRepository.Load()

	if result.IsSuccess() {
		deployment := result.Result.(model.DeploymentEntity)
		s.setDefaultType(&deployment)
		if err := s.storeService.AddDataDeployment(&deployment); err != nil {
			return &model.DeploymentEntity{}, err
		}

		return &deployment, nil
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
