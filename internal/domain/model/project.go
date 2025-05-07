package model

// Project representa la estructura de un proyecto dentro del sistema.
// Contiene información relevante para la identificación y configuración
// del proyecto así como sus dependencias.
type Project struct {
	Organization string                `yaml:"organization"`
	ProjectID    string                `yaml:"projectId"`
	ProjectName  string                `yaml:"projectName"`
	Version      string                `yaml:"version"`
	TeamName     string                `yaml:"teamName"`
	//Dependencies map[string]Dependency `yaml:"dependencies"`
}

// NewProject crea una nueva instancia de Project con valores predeterminados.
// Si organization o teamName están vacíos, se utilizan valores por defecto.
func NewProject(organization, projectID, projectName, teamName string) *Project {
	if organization == "" {
		organization = DefaultOrganization
	}
	if teamName == "" {
		teamName = DefaultTeamName
	}

	return &Project{
		Organization: organization,
		ProjectID:    projectID,
		ProjectName:  projectName,
		Version:      DefaultVersion,
		TeamName:     teamName,
		//Dependencies: make(map[string]Dependency),
	}
}

/* func (p *Project) AddDependency(key string, dependency Dependency) {
	if p.Dependencies != nil {
		p.Dependencies[key] = dependency
	}
} */

func (p *Project) IsComplete() bool {
	return p.Organization != "" && p.ProjectID != "" &&
		p.ProjectName != "" && p.Version != ""
}
