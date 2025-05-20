package service

import (
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/router"
	"deploy/internal/domain/validator"
	"errors"
)

var (
	ErrDeploymentNotFound = errors.New(constant.MsgDeploymentNotFound)
)

type DeploymentServiceInterface interface {
	Load() (*model.Deployment, error)
}

type DeploymentService struct {
	yamlRepository repository.YamlRepository
	fileRepository repository.FileRepository
	router         *router.Router
}

func NewDeploymentService(
	yamlRepository repository.YamlRepository,
	fileRepository repository.FileRepository,
	router *router.Router,
) DeploymentServiceInterface {
	return &DeploymentService{
		yamlRepository: yamlRepository,
		router:         router,
		fileRepository: fileRepository,
	}
}

func (s *DeploymentService) Load() (*model.Deployment, error) {
	path := s.router.GetFullPathDeploymentFile()

	exists := s.fileRepository.ExistsFile(path)
	if !exists {
		return &model.Deployment{}, ErrDeploymentNotFound
	}

	var deployment *model.Deployment
	err := s.yamlRepository.Load(path, &deployment)
	if err != nil {
		return &model.Deployment{}, err
	}

	deployment = s.setDefaultValues(deployment)
	deployment = s.setSetupStep(deployment)

	return deployment, nil
}

func (s *DeploymentService) setSetupStep(deployment *model.Deployment) *model.Deployment {
	if deployment.HasType(validator.TypeContainer) {
		deployment.Steps = append([]model.Step{	
			{
				Name:    validator.TypeSetup,
				Type:    validator.TypeSetup,
				Timeout: "30s",
				Then:    validator.ThenFinish,
			},
		}, deployment.Steps...)
	}
	return deployment
}

func (s *DeploymentService) setDefaultValues(deployment *model.Deployment) *model.Deployment {
	for i := range deployment.Steps {
		if deployment.Steps[i].Type == "" {
			deployment.Steps[i].Type = validator.TypeCommand
		}
	}
	return deployment
}

