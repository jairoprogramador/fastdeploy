package model

import (
	"deploy/internal/domain/constant"
)

// Project representa la estructura de un proyecto dentro del sistema.
// Contiene información relevante para la identificación y configuración
// del proyecto así como sus dependencias.
type Project struct {
	Organization string                `yaml:"organization"`
	ProjectID    string                `yaml:"projectId"`
	Name         string                `yaml:"name"`
	Version      string                `yaml:"version"`
	TeamName     string                `yaml:"teamName"`
	//Dependencies map[string]Dependency `yaml:"dependencies"`
}

// NewProject crea una nueva instancia de Project con valores predeterminados.
// Si organization o teamName están vacíos, se utilizan valores por defecto.
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
		//Dependencies: make(map[string]Dependency),
	}
}

func (p *Project) IsComplete() bool {
	return p.Organization != "" && p.ProjectID != "" &&
		p.Name != "" && p.Version != ""
}
