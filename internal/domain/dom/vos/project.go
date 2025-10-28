package vos

import (
	"errors"
	"time"
)

const (
	DefaultProjectVersion      = "1.0.0"
	DefaultProjectTeam         = "shikigami"
	DefaultProjectDescription  = "mi despliegue con fastdeploy"
	DefaultProjectOrganization = "fastdeploy"
)

type ProjectID string

type Project struct {
	id           ProjectID
	name         string
	version      string
	revision     string
	team         string
	description  string
	organization string
}

func NewProject(id ProjectID, name, version, description, team, organization string) (Project, error) {
	if id == "" {
		return Project{}, errors.New("el ID del proyecto no puede estar vacío")
	}
	if name == "" {
		return Project{}, errors.New("el nombre del proyecto no puede estar vacío")
	}
	if version == "" {
		version = DefaultProjectVersion
	}
	if team == "" {
		team = DefaultProjectTeam
	}
	if organization == "" {
		organization = DefaultProjectOrganization
	}
	if description == "" {
		description = DefaultProjectDescription
	}
	return Project{
		id:           id,
		name:         name,
		version:      version,
		revision:     time.Now().Format("20060102150405"),
		description:  description,
		team:         team,
		organization: organization,
	}, nil
}

func (p Project) ID() ProjectID        { return p.id }
func (p Project) IdString() string     { return string(p.id) }
func (p Project) Name() string         { return p.name }
func (p Project) Description() string  { return p.description }
func (p Project) Team() string         { return p.team }
func (p Project) Version() string      { return p.version }
func (p Project) Revision() string     { return p.revision }
func (p Project) Organization() string { return p.organization }

func (p Project) WithRevision(revision string) Project {
	p.revision = revision
	return p
}
