package service

import (
	"crypto/sha256"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/router"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrProjectNotFound      = errors.New(constant.MsgProjectNotFound)
	ErrProjectCanNotBeNull  = errors.New(constant.MsgProjectCanNotBeNull)
	ErrProjectLoad          = errors.New(constant.MsgProjectLoad)
	ErrProjectNotComplete   = errors.New(constant.MsgProjectNotComplete)
	ErrProjectName          = errors.New(constant.MsgProjectName)
	ErrInvalidProjectData   = errors.New(constant.MsgProjectInvalidData)
	ErrGlobalConfigCreating = errors.New(constant.MsgGlobalConfigCreating)
)

type ProjectServiceInterface interface {
	Initialize() (string, error)
	Load() (*model.Project, error)
	Save(project *model.Project) error
}

type projectService struct {
	yamlRepository      repository.YamlRepository
	globalConfigService GlobalConfigServiceInterface
	fileRepository      repository.FileRepository
	router              *router.Router
	muProjectService    sync.RWMutex
	logStore            *model.LogStore
}

func NewProjectService(
	yamlRepository repository.YamlRepository,
	globalConfigService GlobalConfigServiceInterface,
	fileRepository repository.FileRepository,
	router *router.Router,
	logStore *model.LogStore,
) ProjectServiceInterface {
	return &projectService{
			yamlRepository:      yamlRepository,
			globalConfigService: globalConfigService,
		router:              router,
			fileRepository:      fileRepository,
		logStore:            logStore,
}
}

func (s *projectService) Initialize() (string, error) {
	s.logStore.StartStep("")
	if _, err := s.Load(); err != nil {
		if errors.Is(err, ErrProjectNotFound) || errors.Is(err, ErrProjectNotComplete) {

			project, err := s.createModelProject()
			if err != nil {
				return "", fmt.Errorf("%v", err)
			}

			if err = s.Save(project); err != nil {
				return "", fmt.Errorf("%v", err)
			}
			return constant.MsgInitializeSuccess, nil
		}
		return "", fmt.Errorf("%v", err)
	}
	return constant.MsgInitializeExists, nil
}

func (s *projectService) Load() (*model.Project, error) {
	path := s.router.GetPathProjectFile()

	exists := s.fileRepository.ExistsFile(path)
	if !exists {
		return &model.Project{}, ErrProjectNotFound
	}

	var project model.Project
	err := s.yamlRepository.Load(path, &project)
	if err != nil {
		return &model.Project{}, err
	}

	if !project.IsComplete() {
		return &model.Project{}, ErrProjectNotComplete
	}

	return &project, nil
}

func (s *projectService) Save(project *model.Project) error {
	if project == nil {
		return ErrProjectCanNotBeNull
	}

	if !project.IsComplete() {
		return ErrProjectNotComplete
	}

	path := s.router.GetPathProjectFile()

	if err := s.fileRepository.DeleteFile(path); err != nil {
		return err
	}

	if err := s.yamlRepository.Save(path, project); err != nil {
		return err
	}

	return nil
}

func (s *projectService) createModelProject() (*model.Project, error) {
	projectName, err := s.router.GetProjectName()
	if err != nil || projectName == "" {
		return &model.Project{}, ErrProjectName
	}

	projectId := s.generateProjectId(projectName)

	globalConfig, err := s.getGlobalConfig()
	if err != nil {
		return &model.Project{}, err
	}

	return model.NewProject(globalConfig.Organization, projectId, projectName, globalConfig.TeamName), nil
}

func (s *projectService) getGlobalConfig() (model.GlobalConfig, error) {
	globalConfig, err := s.globalConfigService.Load()

	if err != nil {
		globalConfig = *model.NewGlobalConfigDefault()
		if err = s.globalConfigService.Create(&globalConfig); err != nil {
			return model.GlobalConfig{}, ErrGlobalConfigCreating
		}
	}
	
	return globalConfig, nil
}

func (s *projectService) generateProjectId(prefix string) string {
	timestamp := time.Now().UnixNano()
	data := []byte(fmt.Sprintf("%s-%d", prefix, timestamp))
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
