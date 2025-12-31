package vos

import "errors"

type ProjectData struct {
	name         string
	organization string
	team         string
	description  string
	version      string
}

func NewProjectData(name, organization, team, description, version string) (ProjectData, error) {
	if name == "" {
		return ProjectData{}, errors.New("name is required")
	}
	if organization == "" {
		return ProjectData{}, errors.New("organization is required")
	}
	if team == "" {
		return ProjectData{}, errors.New("team is required")
	}
	return ProjectData{
		name:         name,
		organization: organization,
		team:         team,
		description:  description,
		version:      version,
	}, nil
}

func (p ProjectData) Name() string {
	return p.name
}

func (p ProjectData) Organization() string {
	return p.organization
}

func (p ProjectData) Team() string {
	return p.team
}

func (p ProjectData) Description() string {
	return p.description
}

func (p ProjectData) Version() string {
	return p.version
}
