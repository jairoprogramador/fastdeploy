package service

import (
	"deploy/internal/domain/constant"
	"deploy/internal/domain/engine/model"
	"deploy/internal/domain/engine/validator"
	"deploy/internal/domain/repository"
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
	router               *PathService
}

func NewDeploymentService(
	deploymentRepository repository.DeploymentRepository,
	router *PathService,
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

	return deployment, nil
}
