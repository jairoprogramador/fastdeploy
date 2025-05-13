package service

import (
	"crypto/sha256"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Errores personalizados del servicio de proyecto
var (
	ErrProjectNotFound      = errors.New(constant.MsgProjectNotFound)
	ErrProjectNotComplete   = errors.New(constant.MsgProjectNotComplete)
	ErrProjectName          = errors.New(constant.MsgProjectName)
	ErrProjectCreating      = errors.New(constant.MsgProjectCreating)
	ErrInvalidProjectData   = errors.New(constant.MsgProjectInvalidData)
	ErrGlobalConfigCreating = errors.New(constant.MsgGlobalConfigCreating)
)

type ProjectServiceInterface interface {
	Initialize() *model.Response
	Load() (model.Project, error)
	Create() (model.Project, error)
	SetProjectRepository(projectRepo repository.ProjectRepository)
}

// ProjectService maneja la lógica de negocio relacionada con los proyectos
type ProjectService struct {
	projectRepo         repository.ProjectRepository
	globalConfigService GlobalConfigServiceInterface
	muProjectService    sync.RWMutex
}

var (
	instanceProjectService     *ProjectService
	instanceOnceProjectService sync.Once
)

// GetInstance retorna la instancia única del ProjectService
func GetProjectService(projectRepo repository.ProjectRepository,
	globalConfigService GlobalConfigServiceInterface) ProjectServiceInterface {
	instanceOnceProjectService.Do(func() {
		instanceProjectService = &ProjectService{
			projectRepo:         projectRepo,
			globalConfigService: globalConfigService,
		}
	})
	return instanceProjectService
}

// SetProjectRepository establece el repositorio de proyectos
func (s *ProjectService) SetProjectRepository(projectRepo repository.ProjectRepository) {
	s.muProjectService.Lock()
	defer s.muProjectService.Unlock()
	s.projectRepo = projectRepo
}

// Initialize inicializa el proyecto, creándolo si no existe
func (s *ProjectService) Initialize() *model.Response {
	project, err := s.Load()
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) || errors.Is(err, ErrProjectNotComplete) {
			project, err = s.Create()
			if err != nil {
				return model.GetNewResponseError(err)
			}
			if !project.IsComplete() {
				return model.GetNewResponseMessage(constant.MsgProjectNotComplete)
			}
			return model.GetNewResponseMessage(constant.MsgInitializeSuccess)
		}
		return model.GetNewResponseError(err)
	}
	return model.GetNewResponseMessage(constant.MsgInitializeExists)
}

// Load carga el proyecto desde el repositorio
func (s *ProjectService) Load() (model.Project, error) {
	project, err := s.projectRepo.Load()
	if err != nil {
		return model.Project{}, ErrProjectNotFound
	}

	if err := s.validateProject(&project); err != nil {
		return model.Project{}, err
	}

	return project, nil
}

// Create crea un nuevo proyecto
func (s *ProjectService) Create() (model.Project, error) {
	globalConfig, err := s.loadOrCreateGlobalConfig()
	if err != nil {
		return model.Project{}, err
	}

	projectName, err := s.projectRepo.GetProjectName()
	if err != nil {
		return model.Project{}, ErrProjectName
	}

	if projectName == "" {
		return model.Project{}, ErrProjectName
	}

	projectId := s.generateProjectId(projectName)

	project := model.NewProject(globalConfig.Organization, projectId, projectName, globalConfig.TeamName)

	if err := s.validateProject(project); err != nil {
		return model.Project{}, err
	}

	if err := s.projectRepo.Create(project); err != nil {
		return model.Project{}, ErrProjectCreating
	}

	return *project, nil
}

// validateProject valida que el proyecto tenga todos los campos requeridos
func (s *ProjectService) validateProject(project *model.Project) error {
	if project == nil {
		return ErrInvalidProjectData
	}
	if !project.IsComplete() {
		return ErrInvalidProjectData
	}
	return nil
}

// loadOrCreateGlobalConfig carga o crea la configuración global
func (s *ProjectService) loadOrCreateGlobalConfig() (model.GlobalConfig, error) {
	globalConfig, err := s.globalConfigService.Load()
	if err != nil {
		globalConfig = *model.NewGlobalConfigDefault()
		if _, err = s.globalConfigService.Create(&globalConfig); err != nil {
			return model.GlobalConfig{}, ErrGlobalConfigCreating
		}
	}
	return globalConfig, nil
}

// generateProjectId genera un ID único para el proyecto
func (s *ProjectService) generateProjectId(prefix string) string {
	timestamp := time.Now().UnixNano()
	data := []byte(fmt.Sprintf("%s-%d", prefix, timestamp))
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
