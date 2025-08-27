package entity

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/values"
)

type Configuration struct {
	nameOrganization values.NameOrganization
	team         values.Team
	repository   values.Repository
	technology   values.NameTechnology
}

func NewConfiguration(
	nameOrganization values.NameOrganization,
	team values.Team,
	deployRepo values.Repository,
	technology values.NameTechnology,
) Configuration {
	return Configuration{
		nameOrganization: nameOrganization,
		team:         team,
		repository:   deployRepo,
		technology:   technology,
	}
}

func NewDefaultConfiguration() Configuration {
	org := values.NewDefaultNameOrganization()
	team := values.NewDefaultTeam()
	repo := values.NewDefaultRepository()
	tech := values.NewDefaultNameTechnology()
	return NewConfiguration(org, team, repo, tech)
}

func (c Configuration) GetNameOrganization() values.NameOrganization {
	return c.nameOrganization
}

func (c Configuration) GetTeam() values.Team {
	return c.team
}

func (c Configuration) GetRepository() values.Repository {
	return c.repository
}

func (c Configuration) GetTechnology() values.NameTechnology {
	return c.technology
}
