package application

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/vos"
)

type ProjectService struct {
	projectRepo ports.ProjectRepository
	gitCloner   ports.ClonerTemplate
}

func NewProjectService(
	projectRepo ports.ProjectRepository,
	gitCloner ports.ClonerTemplate) *ProjectService {

	return &ProjectService{
		projectRepo: projectRepo,
		gitCloner:   gitCloner,
	}
}

func (s *ProjectService) Initialize(
	ctx context.Context, projectLocalPath, repositoriesLocalPath string) (*aggregates.Project, error) {
	projectConfigPath := filepath.Join(projectLocalPath, "fdconfig.yaml")

	projectDTO, err := s.projectRepo.Load(ctx, projectConfigPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo cargar la configuración del proyecto: %w", err)
	}

	projectData, err := vos.NewProjectData(
		projectDTO.Name, projectDTO.Organization,
		projectDTO.Team, projectDTO.Description, projectDTO.Version)

	if err != nil {
		return nil, fmt.Errorf("datos del proyecto inválidos: %w", err)
	}
	templateRepo, err := vos.NewTemplateRepository(projectDTO.TemplateURL, projectDTO.TemplateRef)
	if err != nil {
		return nil, fmt.Errorf("datos del repositorio de plantillas inválidos: %w", err)
	}
	projectID := vos.NewProjectID(projectDTO.ID)

	project := aggregates.NewProject(projectID, projectData, templateRepo, projectLocalPath, repositoriesLocalPath)

	if project.SyncID() {
		fmt.Println("El ID del proyecto ha cambiado. Actualizando fdconfig.yaml...")
		projectDTO.ID = project.ID().String()
		if err := s.projectRepo.Save(ctx, projectConfigPath, projectDTO); err != nil {
			return nil, fmt.Errorf("no se pudo guardar el ID del proyecto actualizado: %w", err)
		}
	}

	err = s.gitCloner.EnsureCloned(ctx, project.TemplateRepo().URL(),
		project.TemplateRepo().Ref(), project.TemplateLocalPath())

	if err != nil {
		return nil, fmt.Errorf("no se pudo clonar el repositorio de plantillas: %w", err)
	}

	return project, nil
}
