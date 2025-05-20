package model

type Dependency struct {
	Organization string   `yaml:"organization"`
	ProjectID    string   `yaml:"projectId"`
	ProjectName  string   `yaml:"projectName"`
	Version      string   `yaml:"version"`
	TeamName     string   `yaml:"teamName"`
	Required 	 bool     `yaml:"required"`
}

func NewDependency(organization, projectID, projectName, version, teamName string) *Dependency {
	return &Dependency{
		Organization: organization,
		ProjectID:    projectID,
		ProjectName:  projectName,
		Version:      version,
		TeamName:     teamName,
		Required:     true,
	}
}
