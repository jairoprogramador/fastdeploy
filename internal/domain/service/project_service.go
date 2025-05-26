package service

import (
	"context"
	"crypto/sha256"
	"github.com/jairoprogramador/fastdeploy/internal/domain/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/repository"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrProjectCanNotBeNull = errors.New(constant.ErrorProjectCanNotBeNull)
	ErrProjectNotComplete  = errors.New(constant.ErrorProjectNotComplete)
)

type ProjectService interface {
	Initialize() model.DomainResultEntity
	Start() model.DomainResultEntity

	Load() (*model.ProjectEntity, error)
	Save(project *model.ProjectEntity) error
	GetFullPathResource() (string, error)
}

type projectService struct {
	engine            *engine.Engine
	deploymentService DeploymentService
	projectRepository repository.ProjectRepository
	configService     ConfigService
	router            port.PathService
}

func NewProjectService(
	projectRepository repository.ProjectRepository,
	deploymentService DeploymentService,
	engine *engine.Engine,
	configService ConfigService,
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

func (s *projectService) Start() model.DomainResultEntity {
	if projectEntity, err := s.Load(); err == nil {

		deploymentEntity, err := s.deploymentService.Load()
		if err != nil {
			return model.NewErrorApp(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := s.engine.Execute(ctx, deploymentEntity, projectEntity); err != nil {
			return model.NewErrorApp(err)
		}

		return model.NewResultApp(constant.SuccessInitializeProject)
	}
	return model.NewResultApp(constant.ErrorProjectLoad)
}

func (s *projectService) Initialize() model.DomainResultEntity {
	if _, err := s.Load(); err != nil {
		if err != nil {
			return model.NewErrorApp(err)
		}

		projectEntity, err := s.newProjectEntity()
		if err != nil {
			return model.NewErrorApp(err)
		}

		if err = s.Save(projectEntity); err != nil {
			return model.NewErrorApp(err)
		}
		return model.NewResultApp(constant.SuccessInitializeProject)
	}
	return model.NewResultApp(constant.MsgInitializeExists)
}

func (s *projectService) Load() (*model.ProjectEntity, error) {
	result := s.projectRepository.Load()
	if result.IsSuccess() {
		project := result.Result.(*model.ProjectEntity)
		if !project.IsComplete() {
			return &model.ProjectEntity{}, ErrProjectNotComplete
		}
		return project, nil
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

func (s *projectService) GetFullPathResource() (string, error) {
	result := s.projectRepository.GetFullPathResource()
	if !result.IsSuccess() {
		return "", result.Error
	}
	return result.Result.(string), result.Error
}

func (s *projectService) newProjectEntity() (*model.ProjectEntity, error) {
	result := s.projectRepository.GetName()
	if !result.IsSuccess() {
		return &model.ProjectEntity{}, result.Error
	}

	projectName := result.Result.(string)
	projectId := s.GenerateID(projectName)

	configEntity, err := s.getConfigEntity()
	if err != nil {
		return &model.ProjectEntity{}, err
	}

	projectEntity := model.NewProjectEntity(projectId, projectName)
	projectEntity.Organization = configEntity.Organization
	projectEntity.TeamName = configEntity.TeamName

	return projectEntity, nil
}

func (s *projectService) getConfigEntity() (*model.ConfigEntity, error) {
	configEntity, err := s.configService.Load()

	if err != nil {
		configEntity, err = s.newConfigEntity()
		if err != nil {
			return &model.ConfigEntity{}, err
		}
	}

	if !configEntity.IsComplete() {
		configEntity, err = s.newConfigEntity()
		if err != nil {
			return &model.ConfigEntity{}, err
		}
	}

	return configEntity, nil
}

func (s *projectService) newConfigEntity() (*model.ConfigEntity, error) {
	configEntity := model.NewConfigEntity()
	if err := s.configService.Save(configEntity); err != nil {
		return &model.ConfigEntity{}, err
	}
	return configEntity, nil
}

func (s *projectService) GenerateID(prefix string) string {
	timestamp := time.Now().UnixNano()
	data := []byte(fmt.Sprintf("%s-%d", prefix, timestamp))
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
