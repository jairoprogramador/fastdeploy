package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/technology"
)

type Project struct {
	id           ProjectID
	name         ProjectName
	organization Organization
	team         Team
	repository   repository.Repository
	technology   technology.Technology
	deployment   deployment.Deployment
}

func NewProject(
	id ProjectID,
	name ProjectName,
	organization Organization,
	team Team,
	repository repository.Repository,
	tech technology.Technology,
	deployment deployment.Deployment,
) Project {
	return Project{
		id:           id,
		name:         name,
		organization: organization,
		team:         team,
		repository:   repository,
		technology:   tech,
		deployment:   deployment,
	}
}

func (p Project) GetID() ProjectID {
	return p.id
}

func (p Project) GetName() ProjectName {
	return p.name
}

func (p Project) GetOrganization() Organization {
	return p.organization
}

func (p Project) GetTeam() Team {
	return p.team
}

func (p Project) GetRepository() repository.Repository {
	return p.repository
}

func (p Project) GetTechnology() technology.Technology {
	return p.technology
}

func (p Project) GetDeployment() deployment.Deployment {
	return p.deployment
}

func (p Project) GetFullName() string {
	return p.organization.Value() + "/" + p.name.Value()
}

func (p Project) IncrementDeploymentVersion() Project {
	newDeployment := p.deployment.IncrementVersion()
	return Project{
		id:           p.id,
		name:         p.name,
		organization: p.organization,
		team:         p.team,
		repository:   p.repository,
		technology:   p.technology,
		deployment:   newDeployment,
	}
}

func (p Project) IsValid() bool {
	return p.id.Value() != "" &&
		p.name.Value() != "" &&
		p.organization.Value() != "" &&
		p.team.Value() != "" &&
		p.repository.IsValid() &&
		p.technology.IsValid() &&
		p.deployment.IsValid()
}

func (p Project) Equals(other Project) bool {
	return p.id.Value() == other.id.Value() &&
		p.name.StringValueObject.Equals(other.name.StringValueObject) &&
		p.organization.StringValueObject.Equals(other.organization.StringValueObject) &&
		p.team.StringValueObject.Equals(other.team.StringValueObject) &&
		p.repository.Equals(other.repository) &&
		p.technology.Equals(other.technology) &&
		p.deployment.Equals(other.deployment)
}
