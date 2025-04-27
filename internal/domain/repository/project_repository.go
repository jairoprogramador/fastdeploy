package repository

import "deploy/internal/domain/model"

type ProjectRepository interface {
	Save(project *model.Project) *model.Response
	Load() (model.Project, error)
	IsInitialized(RootDirectory, nameProjectFile string) bool
	GetProjectId() (string, error)
	GetTeamName() string
	GetOrganizationName() string
	SaveDockerfileTemplate() *model.Response
	SaveDockercomposeTemplate() *model.Response
}
