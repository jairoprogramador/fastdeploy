package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/dom/dto"
)

func ProjectToDomain(dto dto.ProjectDTO) (vos.Project, error) {
	projectId := vos.ProjectID(dto.ID)
	return vos.NewProject(
		projectId,
		dto.Name,
		dto.Version,
		dto.Description,
		dto.Team,
	)
}

func ProjectToDTO(project vos.Project) dto.ProjectDTO {
	return dto.ProjectDTO{
		ID:          string(project.ID()),
		Name:        project.Name(),
		Version:     project.Version(),
		Description: project.Description(),
		Team:        project.Team(),
	}
}
