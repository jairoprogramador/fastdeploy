package service

import (
	"crypto/sha256"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/variable"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Errores personalizados del servicio de proyecto
var (
	ErrProjectNotFound      = errors.New(constant.MsgProjectNotFound)
	ErrProjectLoad          = errors.New(constant.MsgProjectLoad)
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
	fileRepository      repository.FileRepository
	globalConfigService GlobalConfigServiceInterface
	store               *variable.VariableStore
	muProjectService    sync.RWMutex
}

var (
	instanceProjectService     *ProjectService
	instanceOnceProjectService sync.Once
)

// GetInstance retorna la instancia única del ProjectService
func GetProjectService(projectRepo repository.ProjectRepository,
	fileRepository repository.FileRepository,
	globalConfigService GlobalConfigServiceInterface) ProjectServiceInterface {
	instanceOnceProjectService.Do(func() {
		instanceProjectService = &ProjectService{
			projectRepo:         projectRepo,
			fileRepository:      fileRepository,
			store:               getStoreProject(),
			globalConfigService: globalConfigService,
		}
	})
	return instanceProjectService
}

func getStoreProject() *variable.VariableStore{
	store := variable.GetVariableStore()
	store.AddVariableGlobal(constant.VAR_PROJECT_ROOT_DIRECTORY, constant.ProjectRootDirectory)
	store.AddVariableGlobal(constant.VAR_PROJECT_FILE_NAME, constant.ProjectFileName)
	return store
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
	pathProjectFile := s.fileRepository.GetFullPathProjectFile(s.store)
	exists := s.fileRepository.ExistsFile(pathProjectFile)
	if !exists {
		return model.Project{}, ErrProjectNotFound
	}

	project, err := s.projectRepo.Load(pathProjectFile)
	if err != nil {
		return model.Project{}, ErrProjectLoad
	}

	if !project.IsComplete() {
		return model.Project{}, ErrProjectNotComplete
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
	if err != nil || projectName == "" {
		return model.Project{}, ErrProjectName
	}

	projectId := s.generateProjectId(projectName)

	project := model.NewProject(globalConfig.Organization, projectId, projectName, globalConfig.TeamName)

	if !project.IsComplete(){
		return model.Project{}, ErrProjectNotComplete
	}

	pathProjectFile := s.fileRepository.GetFullPathProjectFile(s.store)
	if err := s.fileRepository.DeleteFile(pathProjectFile); err != nil {
		return model.Project{},err
	}
	if err := s.projectRepo.Create(pathProjectFile, project); err != nil {
		return model.Project{}, ErrProjectCreating
	}

	return *project, nil
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
