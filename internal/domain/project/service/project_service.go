package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	service3 "github.com/jairoprogramador/fastdeploy/internal/domain/config/service"
	service2 "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine"
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
	Start() result.DomainResult
}

type projectService struct {
	projectRepository repository.ProjectRepository
	configService     service3.ConfigService
	engine            *engine.Engine
	deploymentService service2.DeploymentService
}

func NewProjectService(
	projectRepository repository.ProjectRepository,
	deploymentService service2.DeploymentService,
	engine *engine.Engine,
	configService service3.ConfigService,
) ProjectService {
	return &projectService{
		projectRepository: projectRepository,
		configService:     configService,
		deploymentService: deploymentService,
		engine:            engine,
	}
}

func (s *projectService) Start() result.DomainResult {
	if project, err := s.Load(); err == nil {

		deployment, err := s.deploymentService.Load()
		if err != nil {
			return result.NewErrorApp(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := s.engine.Execute(ctx, deployment, project); err != nil {
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

		project, err := s.createProject()
		if err != nil {
			return result.NewErrorApp(err)
		}

		if err = s.Save(project); err != nil {
			return result.NewErrorApp(err)
		}
		return result.NewResultApp(constant.SuccessInitializeProject)
	}
	return result.NewResultApp(constant.SuccessInitializeExists)
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
