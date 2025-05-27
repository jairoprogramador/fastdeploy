package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	entityConfig "github.com/jairoprogramador/fastdeploy/internal/domain/config/entity"
	service3 "github.com/jairoprogramador/fastdeploy/internal/domain/config/service"
	service2 "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
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
	Start() result.DomainResult

	Load() (*entity.ProjectEntity, error)
	Save(project *entity.ProjectEntity) error
	GetFullPathResource() (string, error)
}

type projectService struct {
	engine            *engine.Engine
	deploymentService service2.DeploymentService
	projectRepository repository.ProjectRepository
	configService     service3.ConfigService
	router            port.PathService
}

func NewProjectService(
	projectRepository repository.ProjectRepository,
	deploymentService service2.DeploymentService,
	engine *engine.Engine,
	configService service3.ConfigService,
	router port.PathService,
) ProjectService {
	return &projectService{
		projectRepository: projectRepository,
		configService:     configService,
		router:            router,
		deploymentService: deploymentService,
		engine:            engine,
	}
}

func (s *projectService) Start() result.DomainResult {
	if projectEntity, err := s.Load(); err == nil {

		deploymentEntity, err := s.deploymentService.Load()
		if err != nil {
			return result.NewErrorApp(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := s.engine.Execute(ctx, deploymentEntity, projectEntity); err != nil {
			return result.NewErrorApp(err)
		}

		return result.NewResultApp(constant.SuccessInitializeProject)
	}
	return result.NewResultApp(constant.ErrorProjectLoad)
}

func (s *projectService) Initialize() result.DomainResult {
	if _, err := s.Load(); err != nil {
		if err != nil {
			return result.NewErrorApp(err)
		}

		projectEntity, err := s.newProjectEntity()
		if err != nil {
			return result.NewErrorApp(err)
		}

		if err = s.Save(projectEntity); err != nil {
			return result.NewErrorApp(err)
		}
		return result.NewResultApp(constant.SuccessInitializeProject)
	}
	return result.NewResultApp(constant.SuccessInitializeExists)
}

func (s *projectService) Load() (*entity.ProjectEntity, error) {
	result := s.projectRepository.Load()
	if result.IsSuccess() {
		project := result.Result.(*entity.ProjectEntity)
		if !project.IsComplete() {
			return &entity.ProjectEntity{}, ErrProjectNotComplete
		}
		return project, nil
	}
	return &entity.ProjectEntity{}, result.Error
}

func (s *projectService) Save(projectEntity *entity.ProjectEntity) error {
	if projectEntity == nil {
		return ErrProjectCanNotBeNull
	}

	if !projectEntity.IsComplete() {
		return ErrProjectNotComplete
	}
	return s.projectRepository.Save(projectEntity).Error
}

func (s *projectService) GetFullPathResource() (string, error) {
	result := s.projectRepository.GetFullPathResource()
	if !result.IsSuccess() {
		return "", result.Error
	}
	return result.Result.(string), result.Error
}

func (s *projectService) newProjectEntity() (*entity.ProjectEntity, error) {
	result := s.projectRepository.GetName()
	if !result.IsSuccess() {
		return &entity.ProjectEntity{}, result.Error
	}

	projectName := result.Result.(string)
	projectId := s.GenerateID(projectName)

	configEntity, err := s.getConfigEntity()
	if err != nil {
		return &entity.ProjectEntity{}, err
	}

	projectEntity := entity.NewProjectEntity(projectId, projectName)
	projectEntity.Organization = configEntity.Organization
	projectEntity.TeamName = configEntity.TeamName

	return projectEntity, nil
}

func (s *projectService) getConfigEntity() (*entityConfig.ConfigEntity, error) {
	configEntity, err := s.configService.Load()

	if err != nil {
		configEntity, err = s.newConfigEntity()
		if err != nil {
			return &entityConfig.ConfigEntity{}, err
		}
	}

	if !configEntity.IsComplete() {
		configEntity, err = s.newConfigEntity()
		if err != nil {
			return &entityConfig.ConfigEntity{}, err
		}
	}

	return configEntity, nil
}

func (s *projectService) newConfigEntity() (*entityConfig.ConfigEntity, error) {
	configEntity := entityConfig.NewConfigEntity()
	if err := s.configService.Save(configEntity); err != nil {
		return &entityConfig.ConfigEntity{}, err
	}
	return configEntity, nil
}

func (s *projectService) GenerateID(prefix string) string {
	timestamp := time.Now().UnixNano()
	data := []byte(fmt.Sprintf("%s-%d", prefix, timestamp))
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
