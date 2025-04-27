package model

const teamNameDefault = "Akatsuki"
const organizationNameDefault = "Jailux"

type Project struct {
	ProjectID    string                   `yaml:"projectId"`
	TeamName     string                   `yaml:"teamName"`
	Organization string                   `yaml:"organization"`
	Language     string                   `yaml:"language"`
	Dependencies map[string]Dependency  `yaml:"dependencies"`
	Support      map[string]Support     `yaml:"support"`
}

func GetNewProject(projectId, teamName, organizationName string) *Project {
	if teamName == "" {
		teamName = teamNameDefault
	} 
	if organizationName == "" {
		organizationName = organizationNameDefault
	} 

	return &Project {
		ProjectID:    projectId,
		TeamName:     teamName,
		Organization: organizationName,
		Dependencies: nil,
		Support:      nil,
	}
}