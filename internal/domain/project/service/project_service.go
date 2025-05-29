package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	serviceConfig "github.com/jairoprogramador/fastdeploy/internal/domain/config/service"
	serviceDeploy "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/repository"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
	"time"
)

var (
	ErrProjectCanNotBeNull = errors.New(constant.ErrorProjectCanNotBeNull)
	ErrProjectNotComplete  = errors.New(constant.ErrorProjectNotComplete)
)

type ProjectService interface {
	Initialize() result.DomainResult
	Start(ctx context.Context) result.DomainResult
}

type projectService struct {
	projectRepository repository.ProjectRepository
	configService     serviceConfig.ConfigService
	containerPort     port.ContainerPort
	engine            *engine.Engine
	deploymentService serviceDeploy.DeploymentService
	storeService      service.StoreServicePort
}

func NewProjectService(
	projectRepository repository.ProjectRepository,
	deploymentService serviceDeploy.DeploymentService,
	engine *engine.Engine,
	configService serviceConfig.ConfigService,
	containerPort port.ContainerPort,
	storeService service.StoreServicePort,
) ProjectService {
	return &projectService{
		projectRepository: projectRepository,
		configService:     configService,
		deploymentService: deploymentService,
		engine:            engine,
		containerPort:     containerPort,
		storeService:      storeService,
	}
}

func (s *projectService) Start(ctx context.Context) result.DomainResult {
	if _, err := s.Load(); err != nil {
		return result.NewMessageApp(constant.ErrorProjectLoad)
	}

	deployment, err := s.deploymentService.Load()
	if err != nil {
		return result.NewErrorApp(err)
	}

	exists, err := s.existsContainer(ctx)
	if err != nil {
		return result.NewErrorApp(err)
	}

	if exists {
		if response := s.containerPort.Up(ctx); !response.IsSuccess() {
			return result.NewErrorApp(response.Error)
		}
	} else {
		if err := s.engine.Execute(ctx, deployment); err != nil {
			return result.NewErrorApp(err)
		}
	}

	if response := s.getUrlContainer(ctx); response.IsSuccess() {
		message := fmt.Sprintf(constant.SuccessStartProjectUrl, response.Result.([]string))
		return result.NewMessageApp(message)
	}

	return result.NewMessageApp(constant.SuccessStartProject)
}

func (s *projectService) Initialize() result.DomainResult {
	if _, err := s.Load(); err != nil {
		project, err := s.createProject()
		if err != nil {
			return result.NewErrorApp(err)
		}

		if err = s.Save(project); err != nil {
			return result.NewErrorApp(err)
		}
		return result.NewMessageApp(constant.SuccessInitializeProject)
	}
	return result.NewMessageApp(constant.SuccessInitializeExists)
}

func (s *projectService) Load() (*model.ProjectEntity, error) {
	result := s.projectRepository.Load()
	if result.IsSuccess() {
		project := result.Result.(model.ProjectEntity)
		if !project.IsComplete() {
			return &model.ProjectEntity{}, ErrProjectNotComplete
		}
		if err := s.storeService.AddDataProject(&project); err != nil {
			return &model.ProjectEntity{}, err
		}
		return &project, nil
	}
	return &model.ProjectEntity{}, result.Error
}

func (s *projectService) Save(projectEntity *model.ProjectEntity) error {
	if projectEntity == nil {
		return ErrProjectCanNotBeNull
	}

	if !projectEntity.IsComplete() {
		return ErrProjectNotComplete
	}
	return s.projectRepository.Save(projectEntity).Error
}

func (s *projectService) createProject() (*model.ProjectEntity, error) {
	result := s.projectRepository.GetName()
	if !result.IsSuccess() {
		return &model.ProjectEntity{}, result.Error
	}

	projectName := result.Result.(string)
	projectId := s.generateID(projectName)

	configEntity, err := s.configService.Get()
	if err != nil {
		return &model.ProjectEntity{}, err
	}

	projectEntity := model.NewProjectEntity(projectId, projectName)
	projectEntity.Organization = configEntity.Organization
	projectEntity.TeamName = configEntity.TeamName

	return projectEntity, nil
}

func (s *projectService) generateID(prefix string) string {
	timestamp := time.Now().UnixNano()
	data := []byte(fmt.Sprintf("%s-%d", prefix, timestamp))
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *projectService) getUrlContainer(ctx context.Context) result.InfraResult {
	commitHash := s.storeService.GetStore().Get(constant.KeyCommitHash)
	projectVersion := s.storeService.GetStore().Get(constant.KeyProjectVersion)
	return s.containerPort.GetURLsUp(ctx, commitHash, projectVersion)
}

func (s *projectService) existsContainer(ctx context.Context) (bool, error) {
	commitHash := s.storeService.GetStore().Get(constant.KeyCommitHash)
	projectVersion := s.storeService.GetStore().Get(constant.KeyProjectVersion)

	response := s.containerPort.Exists(ctx, commitHash, projectVersion)
	if response.IsSuccess() {
		return response.Result.(bool), nil
	}
	return false, response.Error
}
