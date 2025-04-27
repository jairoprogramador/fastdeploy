package service

import "deploy/internal/domain"
import "deploy/internal/domain/model"
import "deploy/internal/domain/repository"
import "sync"


type ProjectService struct {
    projectRepo repository.ProjectRepository
}

var (
    instanceProjectService     *ProjectService
    instanceOnceProjectService sync.Once
    mutexProjectService        sync.Mutex
)

func GetProjectService(projectRepo repository.ProjectRepository) *ProjectService {
    instanceOnceProjectService.Do(func() {
        instanceProjectService = &ProjectService{
            projectRepo: projectRepo,
        }
    })
    return instanceProjectService
}

func (s *ProjectService) SetProjectRepository(projectRepo repository.ProjectRepository) {
    mutexProjectService.Lock()
    defer mutexProjectService.Unlock()
    
    s.projectRepo = projectRepo
}

func (s *ProjectService) Exists() bool {
	return s.projectRepo.Exists()
}

func (s *ProjectService) Create() *model.Response {

	if resp := s.createProyect(); resp.Error != nil {
        return resp
    }

	if resp:= s.createDockerfileTemplate(); resp.Error != nil {
        return resp
    }

	if resp:= s.createComposeTemplate(); resp.Error != nil {
        return resp
    }

	return model.GetNewResponseMessage(constants.MessageSuccessInitializingProject)
}

func (s *ProjectService) createProyect() *model.Response {
	projectId, _ := s.projectRepo.GetProjectId()
	teamName := s.projectRepo.GetTeamName()
	organizationName := s.projectRepo.GetOrganizationName()

	project := model.GetNewProject(projectId, teamName, organizationName)
	project.Dependencies = s.getDependencies()
	project.Support = s.getSupports()

	return s.projectRepo.Save(project)
}

func (s *ProjectService) createDockerfileTemplate() *model.Response {
	return s.projectRepo.SaveDockerfileTemplate()
}

func (s *ProjectService) createComposeTemplate() *model.Response {
	return s.projectRepo.SaveDockercomposeTemplate()
}

func (st *ProjectService) getDependencies() map[string]model.Dependency {
	//init crear archivo de dependencia por defecto, ese archivo se le y se crea las dependencias
	ConfigMysql := map[string]string {
		"host":     "localhost",
		"port":     "3306",
		"database": "dev_db",
	}
	dependencyMysql := model.GetNewDependency("database", "mysql", "8.0", true, ConfigMysql)

	ConfigVault := map[string]string {
		"path":     "secrets/dev",
	}
	dependencyVault := model.GetNewDependency("secrets", "vault", "1.0", true, ConfigVault)

	dependencies := map[string]model.Dependency{
		"database": dependencyMysql,
		"secret": dependencyVault,
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
