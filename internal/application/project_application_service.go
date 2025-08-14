package application

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/services"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/repositories"
)

type ProjectApplicationService struct{}

func NewProjectApplicationService() *ProjectApplicationService {
	return &ProjectApplicationService{}
}

// CreateProject es un método para inicializar un nuevo proyecto.
// Debe instanciar los repositorios y servicios de dominio (ProjectService y ConfigService)
// para realizar la operación de inicialización. Debe manejar posibles errores de I/O.
func (pas *ProjectApplicationService) CreateProject(projectName string) error {
	// Instanciar repositorios
	projectRepo := repositories.NewProjectRepository()

	// Instanciar servicios de dominio
	configService := services.NewConfigService()
	projectService := services.NewProjectService()

	// Realizar la operación de inicialización
	_, err := projectService.Initialize(configService, *projectRepo, projectName)
	if err != nil {
		return fmt.Errorf("error creating project: %w", err)
	}

	return nil
}

// GetProject es un método para obtener los datos de un proyecto.
// Debe instanciar el repositorio y usar ProjectService para cargar el proyecto.
func (pas *ProjectApplicationService) GetProject() (project.Project, error) {
	// Instanciar repositorio
	projectRepo := repositories.NewProjectRepository()

	// Instanciar servicio de dominio
	projectService := services.NewProjectService()

	// Cargar el proyecto
	proj, err := projectService.Load(*projectRepo)
	if err != nil {
		return project.Project{}, fmt.Errorf("error loading project: %w", err)
	}

	return proj, nil
}
