package configuration

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/repository"
)

type Configuration struct {
	organization project.Organization
	team         project.Team
	repository   repository.Repository
}

func NewConfiguration(
	organization project.Organization,
	team project.Team,
	deployRepo repository.Repository,
) Configuration {
	return Configuration{
		organization: organization,
		team:         team,
		repository:   deployRepo,
	}
}

func (c Configuration) GetOrganization() project.Organization {
	return c.organization
}

func (c Configuration) GetTeam() project.Team {
	return c.team
}

func (c Configuration) GetRepository() repository.Repository {
	return c.repository
}

func (c Configuration) IsValid() bool {
	return c.organization.Value() != "" &&
		c.team.Value() != "" &&
		c.repository.IsValid()
}

func (c Configuration) Equals(other Configuration) bool {
	return c.organization.StringValueObject.Equals(other.organization.StringValueObject) &&
		c.team.StringValueObject.Equals(other.team.StringValueObject) &&
		c.repository.Equals(other.repository)
}
