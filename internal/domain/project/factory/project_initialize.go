package factory

import (
	config "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entity"
	project "github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/values"
)

type ProjectFactory interface {
	Create(
		configuration config.Configuration,
		idProject string,
		nameProject string,
	) (project.Project, error)
}

type ProjectInitialize struct { }

func NewProjectFactory() ProjectFactory {
	return &ProjectInitialize{}
}

func (pf *ProjectInitialize) Create(
	configuration config.Configuration,
	idProject string,
	nameProject string,
) (project.Project, error) {

	projectId, err := values.NewIdentifier(idProject)
	if err != nil {
		return project.Project{}, err
	}

	projectName, err := values.NewNameProject(nameProject)
	if err != nil {
		return project.Project{}, err
	}

	organization, err := pf.makeNameOrganization(configuration.GetNameOrganization())
	if err != nil {
		return project.Project{}, err
	}

	team, err := pf.makeTeam(configuration.GetTeam())
	if err != nil {
		return project.Project{}, err
	}

	repository, err := pf.makeRepository(configuration.GetRepository())
	if err != nil {
		return project.Project{}, err
	}

	technology, err := pf.makeTechnology(configuration.GetTechnology())
	if err != nil {
		return project.Project{}, err
	}

	return project.NewProject(
		projectId,
		projectName,
		organization,
		team,
		repository,
		technology,
		values.NewDefaultDeployment(),
		values.NewDefaultCategoryProject(),
	), nil
}

func (pf *ProjectInitialize) makeNameOrganization(name values.NameOrganization) (values.NameOrganization, error) {
	if name.IsEmpty() {
		return values.NewDefaultNameOrganization(), nil
	}
	return values.NewNameOrganization(name.Value())
}

func (pf *ProjectInitialize) makeTeam(team values.Team) (values.Team, error) {
	if team.IsEmpty() {
		return values.NewDefaultTeam(), nil
	}

	return values.NewTeam(team.Value())
}

func (pf *ProjectInitialize) makeRepository(repository values.Repository) (values.Repository, error) {
	url := repository.GetURL()
	if url.IsEmpty() {
		url = values.NewDefaultUrlRepository()
	} else {
		var err error
		url, err = values.NewUrlRepository(url.Value())
		if err != nil {
			return values.Repository{}, err
		}
	}

	version := repository.GetVersion()
	if version.IsEmpty() {
		version = values.NewDefaultVersionRepository()
	} else {
		var err error
		version, err = values.NewVersionRepository(version.Value())
		if err != nil {
			return values.Repository{}, err
		}
	}

	return values.NewRepository(url, version), nil
}

func (pf *ProjectInitialize) makeTechnology(technology values.NameTechnology) (values.NameTechnology, error) {
	if technology.IsEmpty() {
		return values.NewDefaultNameTechnology(), nil
	} 
	return values.NewNameTechnology(technology.Value())
}

