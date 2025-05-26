package entity

import (
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

// ProjectEntity represents a project in the system
type ProjectEntity struct {
	Organization string `yaml:"organization"`
	ProjectID    string `yaml:"projectId"`
	Name         string `yaml:"name"`
	Version      string `yaml:"version"`
	TeamName     string `yaml:"teamName"`
}

// NewProjectEntity creates a new ProjectEntity with the given ID and name
func NewProjectEntity(projectID string, name string) *ProjectEntity {
	return &ProjectEntity{
		ProjectID: projectID,
		Name:      name,
		Version:   constant.DefaultVersion,
	}
}

// IsComplete checks if all required fields are set
func (p *ProjectEntity) IsComplete() bool {
	return p.Organization != "" && p.ProjectID != "" &&
		p.Name != "" && p.Version != ""
}