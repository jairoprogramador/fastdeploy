package model

import (
	"deploy/internal/domain/constant"
)

type ProjectEntity struct {
	Organization string `yaml:"organization"`
	ProjectID    string `yaml:"projectId"`
	Name         string `yaml:"name"`
	Version      string `yaml:"version"`
	TeamName     string `yaml:"teamName"`
}

func NewProjectEntity(projectID string, name string) *ProjectEntity {
	return &ProjectEntity{
		ProjectID: projectID,
		Name:      name,
		Version:   constant.DefaultVersion,
	}
}

func (p *ProjectEntity) IsComplete() bool {
	return p.Organization != "" && p.ProjectID != "" &&
		p.Name != "" && p.Version != ""
}
