package mapper

import (
	proEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/entities"
	iProDto "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project/dto"
)

func ProjectToDomain(dto iProDto.ProjectDTO) (*proEnt.Project, error) {
	return proEnt.NewProject(
		proEnt.ProjectID(dto.ID),
		dto.Name,
		dto.Version,
		dto.Description,
		dto.Team,
		dto.Organization)
}

func ProjectToDTO(project *proEnt.Project) iProDto.ProjectDTO {
	return iProDto.ProjectDTO{
		ID: string(project.ID()),
		Name: project.Name(),
		Version: project.Version(),
		Description: project.Description(),
		Team: project.Team(),
		Organization: project.Organization(),
	}
}