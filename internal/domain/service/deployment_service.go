package service

import (
	"deploy/internal/domain/constant"
	"deploy/internal/domain/engine/validator"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/service/router"
	"errors"
)

var (
	ErrDeploymentNotFound = errors.New(constant.MsgDeploymentNotFound)
)

type DeploymentLoader interface {
	Load() (*model.DeploymentEntity, error)
}

type deploymentService struct {
	deploymentRepository repository.DeploymentRepository
	router               *router.Router
}

func NewDeploymentService(
	deploymentRepository repository.DeploymentRepository,
	router *router.Router,
) DeploymentLoader {
	return &deploymentService{
		deploymentRepository: deploymentRepository,
		router:               router,
	}
}

func (s *deploymentService) Load() (*model.DeploymentEntity, error) {
	deployment, err := s.deploymentRepository.Load()
	if err != nil {
		return nil, err
	}

	for i := range deployment.Steps {
		if deployment.Steps[i].Type == "" {
			deployment.Steps[i].Type = validator.TypeCommand
		}
	}

	if deployment.HasType(validator.TypeContainer) {
		setupStep := model.Step{
			Name:    validator.TypeSetup,
			Type:    validator.TypeSetup,
			Timeout: "30s",
			Then:    validator.ThenFinish,
		}
		deployment.Steps = append([]model.Step{setupStep}, deployment.Steps...)
	}

	return deployment, nil
}
