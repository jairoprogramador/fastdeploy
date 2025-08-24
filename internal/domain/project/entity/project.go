package entity

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/values"
)

type Project struct {
	id           values.Identifier
	name         values.NameProject
	organization values.NameOrganization
	team         values.Team
	repository   values.Repository
	technology   values.Technology
	deployment   values.Deployment
}

func NewProject(
	id values.Identifier,
	name values.NameProject,
	organization values.NameOrganization,
	team values.Team,
	repository values.Repository,
	tech values.Technology,
	deployment values.Deployment,
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

func (p Project) GetID() values.Identifier {
	return p.id
}

func (p Project) GetName() values.NameProject {
	return p.name
}

func (p Project) GetOrganization() values.NameOrganization {
	return p.organization
}

func (p Project) GetTeam() values.Team {
	return p.team
}

func (p Project) GetRepository() values.Repository {
	return p.repository
}

func (p Project) GetTechnology() values.Technology {
	return p.technology
}

func (p Project) GetDeployment() values.Deployment {
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

