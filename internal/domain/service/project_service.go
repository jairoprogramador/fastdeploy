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

/* func (s *ProjectService) loadGlobalConfig() (model.GlobalConfig, error) {
	if exists := s.globalConfigRepo.ExistsFile();!exists {
		return model.GlobalConfig{}, errors.New(constants.MsgGlobalConfigNotFound)
	}
	globalConfig, err := s.globalConfigRepo.Load()
	if err != nil {
		return model.GlobalConfig{}, errors.New(constants.MsgGlobalConfigNotFound)
	}
	return globalConfig, nil
} */

/* func (s *ProjectService) createGlobalConfig() error {
	if err := s.globalConfigRepo.RemoveFile(); err != nil {
		return err
	}

	globalConfig := model.NewGlobalConfigDefault()
	return s.globalConfigRepo.Create(globalConfig)

} */

/* func (s *ProjectService) Exists() bool {
	return s.projectRepo.Exists()
} */

/* func (s *ProjectService) Create() *model.Response {

	if resp := s.createProject(); resp.Error != nil {
        return resp
    }

	if resp:= s.createDockerfileTemplate(); resp.Error != nil {
        return resp
    }

	if resp:= s.createComposeTemplate(); resp.Error != nil {
        return resp
    }

	return model.GetNewResponseMessage(constants.MsgInitializeSuccess)
} */

/* func (s *ProjectService) createProjectd() *model.Response {
	projectId, _ := s.projectRepo.GetProjectId()
	projectName, _ := s.projectRepo.GetProjectName()
	teamName := s.projectRepo.GetTeamName()
	organizationName := s.projectRepo.GetOrganizationName()

	project := model.NewProject(projectId, projectName, teamName, organizationName)
	//project.Dependencies = s.getDependencies()
	//project.Support = s.getSupports()

	return s.projectRepo.Save(project)
}

func (s *ProjectService) createDockerfileTemplate() *model.Response {
	return s.projectRepo.SaveDockerfileTemplate()
}

func (s *ProjectService) createComposeTemplate() *model.Response {
	return s.projectRepo.SaveDockercomposeTemplate()
} */

/* func (st *ProjectService) getDependencies() map[string]model.Dependency {

	dependencyMysql := model.NewDependency("Jailux", "1", "mysql", "8.0", "Akatsuki")

	dependencyVault := model.NewDependency("Jailux", "2", "vault", "1.0", "Akatsuki")

	dependencies := map[string]model.Dependency{
		"database": *dependencyMysql,
		"secret": *dependencyVault,
	}
	return dependencies
}

func (st *ProjectService) getSupports() map[string]model.Support {
	//init crear archivo de supports por defecto, ese archivo se le y se crea los supports
	ConfigSonarqube := map[string]string {
		"projectKey":  "project-dev",
	}
	dependencySonarqube := model.GetNewSupport("quality", "sonarqube", "9.0", "http://localhost:9000", ConfigSonarqube)

	supports := map[string]model.Support{
		"quality": dependencySonarqube,
	}
	return supports
}
*/
