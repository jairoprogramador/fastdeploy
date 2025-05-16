package service


import (
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/router"
	"errors"
	"sync"
)

var (
	ErrDeploymentNotFound      = errors.New(constant.MsgDeploymentNotFound)
)

type DeploymentServiceInterface interface {
	Load() (model.Deployment, error)
}

type DeploymentService struct {
	yamlRepository      repository.YamlRepository
	fileRepository      repository.FileRepository
	router              *router.Router
	muDeploymentService    sync.RWMutex
}

var (
	instanceDeploymentService     *DeploymentService
	instanceOnceDeploymentService sync.Once
)

func GetDeploymentService(
	yamlRepository repository.YamlRepository,
	fileRepository repository.FileRepository,) DeploymentServiceInterface {

	instanceOnceDeploymentService.Do(func() {
		instanceDeploymentService = &DeploymentService{
			yamlRepository:      yamlRepository,
			router:              router.GetRouter(),
			fileRepository:      fileRepository,
		}
	})
	return instanceDeploymentService
}

func (s *DeploymentService) SetYamlRepository(yamlRepository repository.YamlRepository) {
	s.muDeploymentService.Lock()
	defer s.muDeploymentService.Unlock()
	s.yamlRepository = yamlRepository
}

func (s *DeploymentService) Load() (model.Deployment, error) {
	path := s.router.GetFullPathDeploymentFile()

	exists := s.fileRepository.ExistsFile(path)
	if !exists {
		return model.Deployment{}, ErrDeploymentNotFound
	}

	var deployment model.Deployment
	err := s.yamlRepository.Load(path, &deployment)
	if err != nil {
		return model.Deployment{}, err
	}

	return deployment, nil
}
