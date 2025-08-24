package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/values"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/dto"
)

func ToDto(project entity.Project) (dto.ProjectDto, error) {
	projectInfo := dto.ProjectInfo{
		ID:   project.GetID().Value(),
		Name: project.GetName().Value(),
	}
	repository := dto.RepositoryInfo{
		URL: project.GetRepository().GetURL().Value(),
	}
	technology := dto.TechnologyInfo{
		Name:    project.GetTechnology().GetName().Value(),
		Version: project.GetTechnology().GetVersion().Value(),
	}
	deployment := dto.DeploymentInfo{
		Version: project.GetDeployment().GetVersion().Value(),
	}

	return dto.ProjectDto{
		Organization: project.GetOrganization().Value(),
		Team:         project.GetTeam().Value(),
		Project:      projectInfo,
		Repository:   repository,
		Technology:   technology,
		Deployment:   deployment,
	}, nil
}

func ToDomain(dto dto.ProjectDto) (entity.Project, error) {

	id, err := values.NewIdentifier(dto.Project.ID)
	if err != nil {
		return entity.Project{}, err
	}

	name, err := values.NewNameProject(dto.Project.Name)
	if err != nil {
		return entity.Project{}, err
	}

	organization, err := values.NewNameOrganization(dto.Organization)
	if err != nil {
		return entity.Project{}, err
	}

	team, err := values.NewTeam(dto.Team)
	if err != nil {
		return entity.Project{}, err
	}

	repositoryUrl, err := values.NewUrlRepository(dto.Repository.URL)
	if err != nil {
		return entity.Project{}, err
	}

	repository := values.NewRepository(repositoryUrl)

	technologyName, err := values.NewNameTechnology(dto.Technology.Name)
	if err != nil {
		return entity.Project{}, err
	}

	technologyVersion, err := values.NewVersionTechnology(dto.Technology.Version)
	if err != nil {
			return entity.Project{}, err
	}

	technology := values.NewTechnology(technologyName, technologyVersion)

	deploymentVersion, err := values.NewVersionDeployment(dto.Deployment.Version)
	if err != nil {
		return entity.Project{}, err
	}

	deployment := values.NewDeployment(deploymentVersion)

	return entity.NewProject(
		id,
		name,
		organization,
		team,
		repository,
		technology,
		deployment,
	), nil
}
