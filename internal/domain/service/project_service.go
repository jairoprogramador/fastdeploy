package service

import (
	"crypto/sha256"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/repository"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrProjectNotFound     = errors.New(constant.MsgProjectNotFound)
	ErrProjectCanNotBeNull = errors.New(constant.MsgProjectCanNotBeNull)
	ErrProjectNotComplete  = errors.New(constant.MsgProjectNotComplete)
)

type ProjectLoader interface {
	Load() (*model.ProjectEntity, error)
}

type ProjectPersister interface {
	Save(project *model.ProjectEntity) error
}

type ProjectInitializer interface {
	Initialize() (string, error)
}

type ProjectService interface {
	ProjectLoader
	ProjectPersister
	ProjectInitializer
	GetFullPathResource() (string, error)
}

type projectService struct {
	logger            *logger.Logger
	projectRepository repository.ProjectRepository
	configService     ConfigService
	router            *PathService
}

func NewProjectService(
	logger *logger.Logger,
	projectRepository repository.ProjectRepository,
	configService ConfigService,
	router *PathService,
) ProjectService {
	return &projectService{
		projectRepository: projectRepository,
		configService:     configService,
		router:            router,
		logger:            logger,
	}
}

func (s *projectService) Initialize() (string, error) {
	if _, err := s.Load(); err != nil {
		if errors.Is(err, ErrProjectNotFound) || errors.Is(err, ErrProjectNotComplete) {
			projectEntity, err := s.newProjectEntity()
			if err != nil {
				return "", err
			}

			if err = s.Save(projectEntity); err != nil {
				return "", err
			}
			return constant.MsgInitializeSuccess, nil
		}
		return "", err
	}
	return constant.MsgInitializeExists, nil
}

func (s *projectService) Load() (*model.ProjectEntity, error) {
	project, err := s.projectRepository.Load()
	if err == nil && project != nil {
		if !project.IsComplete() {
			return &model.ProjectEntity{}, ErrProjectNotComplete
		}
	}
	return project, err
}

func (s *projectService) Save(projectEntity *model.ProjectEntity) error {
	if projectEntity == nil {
		return ErrProjectCanNotBeNull
	}

	if !projectEntity.IsComplete() {
		return ErrProjectNotComplete
	}
	return s.projectRepository.Save(projectEntity)
}

func (s *projectService) GetFullPathResource() (string, error) {
	return s.projectRepository.GetFullPathResource()
}

func (s *projectService) newProjectEntity() (*model.ProjectEntity, error) {
	projectName, err := s.projectRepository.GetName()
	if err != nil || projectName == "" {
		return &model.ProjectEntity{}, err
	}

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
		if errors.Is(err, ErrConfigNotFound) || errors.Is(err, ErrConfigNotComplete) {
			configEntity, err = s.newConfigEntity()
			if err != nil {
				return &model.ConfigEntity{}, err
			}
			return configEntity, nil
		} else {
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
