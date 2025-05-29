package service

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
	"time"
)

type DeploymentService interface {
	Load() (*model.DeploymentEntity, error)
}

type deploymentService struct {
	deploymentRepository repository.DeploymentRepository
	containerPort        port.ContainerPort
	storeService         service.StoreServicePort
}

func NewDeploymentService(
	deploymentRepository repository.DeploymentRepository,
	containerPort port.ContainerPort,
	storeService service.StoreServicePort,
) DeploymentService {
	return &deploymentService{
		deploymentRepository: deploymentRepository,
		containerPort:        containerPort,
		storeService:         storeService,
	}
}

func (s *deploymentService) Load() (*model.DeploymentEntity, error) {
	result := s.deploymentRepository.Load()

	if result.IsSuccess() {
		deployment := result.Result.(model.DeploymentEntity)
		s.setDefaultType(&deployment)
		if err := s.setCheckContainer(&deployment); err != nil {
			return &model.DeploymentEntity{}, err
		}

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

func (s *deploymentService) setCheckContainer(deployment *model.DeploymentEntity) error {
	resultFileCompose := s.containerPort.ExistsFileCompose()
	if !resultFileCompose.IsSuccess() {
		return resultFileCompose.Error
	}
	fileComposeExists := resultFileCompose.Result.(bool)

	containerExists, err := s.existsContainer()
	if err != nil {
		return err
	}

	if fileComposeExists && containerExists {
		setupStep := model.Step{
			Name:    string(model.Check),
			Type:    string(model.Check),
			Timeout: "30s",
			Then:    validator.ThenFinish,
		}
		deployment.Steps = append([]model.Step{setupStep}, deployment.Steps...)
	}
	return nil
}

func (s *deploymentService) existsContainer() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	commitHash := s.storeService.GetStore().Get(constant.KeyCommitHash)
	projectVersion := s.storeService.GetStore().Get(constant.KeyProjectVersion)

	response := s.containerPort.Exists(ctx, commitHash, projectVersion)
	if response.IsSuccess() {
		return response.Result.(bool), nil
	}
	return false, response.Error
}
