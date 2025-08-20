package factories

import (
	config "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
	project "github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/values"
)

type ProjectFactory interface {
	Create(configuration config.Configuration, idProject string, nameProject string) (project.Project, error)
}

type ProjectInitialize struct { }

func NewProjectFactory() ProjectFactory {
	return &ProjectInitialize{}
}

func (pf *ProjectInitialize) Create(configuration config.Configuration, idProject string, nameProject string) (project.Project, error) {

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

	return values.NewRepository(url), nil
}

func (pf *ProjectInitialize) makeTechnology(technology values.Technology) (values.Technology, error) {
	name := technology.GetName()
	if name.IsEmpty() {
		name = values.NewDefaultNameTechnology()
	} else {
		var err error
		name, err = values.NewNameTechnology(name.Value())
		if err != nil {
			return values.Technology{}, err
		}
	}

	version := technology.GetVersion()
	if version.IsEmpty() {
		version = values.NewDefaultVersionTechnology()
	} else {
		var err error
		version, err = values.NewVersionTechnology(version.Value())
		if err != nil {
			return values.Technology{}, err
		}
	}

	return values.NewTechnology(name, version), nil
}

