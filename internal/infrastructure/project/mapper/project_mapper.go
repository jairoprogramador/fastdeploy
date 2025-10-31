package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project/dto"
)

func ProjectToDomain(dto dto.ProjectDTO) (vos.Project, error) {
	return vos.NewProject(
		vos.ProjectID(dto.ID),
		dto.Name,
		dto.Version,
		dto.Description,
		dto.Team,
		dto.Organization)
}

func ProjectToDTO(project vos.Project) dto.ProjectDTO {
	return dto.ProjectDTO{
		ID: string(project.ID()),
		Name: project.Name(),
		Version: project.Version(),
		Description: project.Description(),
		Team: project.Team(),
		Organization: project.Organization(),
	}
}