package model

import (
	"deploy/internal/domain/constant"
)

type Project struct {
	Organization string                `yaml:"organization"`
	ProjectID    string                `yaml:"projectId"`
	Name         string                `yaml:"name"`
	Version      string                `yaml:"version"`
	TeamName     string                `yaml:"teamName"`
}

func NewProject(organization, projectID, name, teamName string) *Project {
	if organization == "" {
		organization = constant.DefaultOrganization
	}
	if teamName == "" {
		teamName = constant.DefaultTeamName
	}

	return &Project{
		Organization: organization,
		ProjectID:    projectID,
		Name:  		  name,
		Version:      constant.DefaultVersion,
		TeamName:     teamName,
	}
}

func (p *Project) IsComplete() bool {
	return p.Organization != "" && p.ProjectID != "" &&
		p.Name != "" && p.Version != ""
}
