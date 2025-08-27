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
		Technology: project.GetTechnology().Value(),
	}

	repository := dto.RepositoryInfo{
		URL: project.GetRepository().GetURL().Value(),
		Version: project.GetRepository().GetVersion().Value(),
	}

	deployment := dto.DeploymentInfo{
		Version: project.GetDeployment().GetVersion().Value(),
	}

	return dto.ProjectDto{
		Organization: project.GetOrganization().Value(),
		Team:         project.GetTeam().Value(),
		Project:      projectInfo,
		Repository:   repository,
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

	technology, err := values.NewNameTechnology(dto.Project.Technology)
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

	repository := values.NewRepository(repositoryUrl, values.NewDefaultVersionRepository())

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
