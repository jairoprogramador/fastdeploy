package services

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/technology"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/repositories"
)

type ProjectService struct{}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// Initialize es el método principal. Debe usar config_service para obtener la configuración,
// crear una nueva instancia de Project con los datos obtenidos, y luego usar project_repo
// para guardar la nueva entidad.
func (ps *ProjectService) Initialize(configService *ConfigService, projectRepo repositories.ProjectRepository, projectName string) (project.Project, error) {
	// Obtener configuración
	configRepo := repositories.NewConfigRepository()
	config, err := configService.Load(*configRepo)
	if err != nil {
		return project.Project{}, err
	}

	// Crear entidades del proyecto
	projectNameEntity, err := project.NewProjectName(projectName)
	if err != nil {
		return project.Project{}, err
	}

	projectID := project.GenerateProjectID(projectName, config.GetOrganization().Value())

	// Crear tecnología por defecto
	techName, _ := technology.NewTechnologyName("default")
	techVersion, _ := technology.NewTechnologyVersion("1.0.0")
	tech := technology.NewTechnology(techName, techVersion)

	// Crear deployment por defecto
	deploymentVersion, _ := deployment.NewDeploymentVersion("1.0.0")
	deployment := deployment.NewDeployment(deploymentVersion)

	// Crear el proyecto
	newProject := project.NewProject(
		projectID,
		projectNameEntity,
		config.GetOrganization(),
		config.GetTeam(),
		config.GetRepository(),
		tech,
		deployment,
	)

	// Guardar el proyecto
	err = ps.saveProject(projectRepo, newProject)
	if err != nil {
		return project.Project{}, err
	}

	return newProject, nil
}

// Load carga la entidad Project desde un archivo usando el repositorio.
func (ps *ProjectService) Load(repo repositories.ProjectRepository) (project.Project, error) {
	_, err := repo.Load()
	if err != nil {
		return project.Project{}, err
	}

	// Aquí se convertiría el diccionario a la entidad Project
	// Por ahora retornamos un proyecto vacío
	return project.Project{}, nil
}

// saveProject guarda el proyecto en el repositorio
func (ps *ProjectService) saveProject(repo repositories.ProjectRepository, proj project.Project) error {
	// Aquí se convertiría la entidad Project a un diccionario
	// Por ahora guardamos un diccionario vacío
	data := make(map[string]interface{})
	return repo.Save(data)
}
