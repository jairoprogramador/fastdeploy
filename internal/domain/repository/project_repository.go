package repository

import "deploy/internal/domain/model"

type ProjectRepository interface {
	Exists() bool
	Save(project *model.Project) *model.Response
	Load() (model.Project, error)
	GetProjectId() (string, error)
	GetTeamName() string
	GetOrganizationName() string
	SaveDockerfileTemplate() *model.Response
	SaveDockercomposeTemplate() *model.Response
}
