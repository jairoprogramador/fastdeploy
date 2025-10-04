package vos

import "errors"

// ProjectID es un hash que identifica unívocamente un proyecto.
type ProjectID string

// Project agrupa los datos que definen un proyecto. Es un Objeto de Valor.
type Project struct {
	id           ProjectID
	name         string
	version      string
	description  string
	team         string
}

func NewProject(id ProjectID, name, version, description, team string) (*Project, error) {
	if id == "" {
		return nil, errors.New("el ID del proyecto no puede estar vacío")
	}
	if name == "" {
		return nil, errors.New("el nombre del proyecto no puede estar vacío")
	}
	if version == "" {
		return nil, errors.New("la versión del proyecto no puede estar vacía")
	}
	if team == "" {
		return nil, errors.New("el equipo del proyecto no puede estar vacío")
	}
	return &Project{
		id:           id,
		name:         name,
		version:      version,
		description:  description,
		team:         team,
	}, nil
}

// Getters para todos los campos...
func (p *Project) ID() ProjectID        { return p.id }
func (p *Project) IdString() string     { return string(p.id) }
func (p *Project) Name() string         { return p.name }
func (p *Project) Description() string  { return p.description }
func (p *Project) Team() string         { return p.team }
func (p *Project) Version() string      { return p.version }