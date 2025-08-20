package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/values"
)

type ProjectInfo struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type TechnologyInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type RepositoryInfo struct {
	URL string `yaml:"url"`
}

type DeploymentInfo struct {
	Version string `yaml:"version"`
}


type ProjectDTO struct {
	Organization string `yaml:"organization"`
	Team     string `yaml:"team"`
	Project  ProjectInfo `yaml:"project"`
	Repository   RepositoryInfo `yaml:"repository"`
	Technology   TechnologyInfo `yaml:"technology"`
	Deployment   DeploymentInfo `yaml:"deployment"`
	
}

func (dto *ProjectDTO) FromDomain(project entities.Project) {
	dto.Organization = project.GetOrganization().Value()

	dto.Project = ProjectInfo{
		ID: project.GetID().Value(),
		Name: project.GetName().Value(),
	}
	dto.Repository = RepositoryInfo{
		URL: project.GetRepository().GetURL().Value(),
	}
	dto.Technology = TechnologyInfo{
		Name: project.GetTechnology().GetName().Value(),
		Version: project.GetTechnology().GetVersion().Value(),
	}
	dto.Deployment = DeploymentInfo{
		Version: project.GetDeployment().GetVersion().Value(),
	}
	dto.Team = project.GetTeam().Value()
}

func (dto ProjectDTO) ToDomain() (entities.Project, error) {

	id, err := values.NewIdentifier(dto.Project.ID)
	if err != nil {
		return entities.Project{}, err
	}

	name, err := values.NewNameProject(dto.Project.Name)
	if err != nil {
		return entities.Project{}, err
	}

	organization, err := values.NewNameOrganization(dto.Organization)
	if err != nil {
		return entities.Project{}, err
	}

	team, err := values.NewTeam(dto.Team)
	if err != nil {
		return entities.Project{}, err
	}

	repositoryUrl, err := values.NewUrlRepository(dto.Repository.URL)
	if err != nil {
		return entities.Project{}, err
	}

	repository := values.NewRepository(repositoryUrl)

	technologyName, err := values.NewNameTechnology(dto.Technology.Name)
	if err != nil {
		return entities.Project{}, err
	}

	technologyVersion, err := values.NewVersionTechnology(dto.Technology.Version)
	if err != nil {
		return entities.Project{}, err
	}

	technology := values.NewTechnology(technologyName, technologyVersion)

	deploymentVersion, err := values.NewVersionDeployment(dto.Deployment.Version)
	if err != nil {
		return entities.Project{}, err
	}

	deployment := values.NewDeployment(deploymentVersion)

	return entities.NewProject(
		id,
		name,
		organization,
		team,
		repository,
		technology,
		deployment,
	), nil
}