package model

type Project struct {
	ProjectId    string                   `yaml:"projectId"`
	TeamName     string                   `yaml:"teamName"`
	Organization string                   `yaml:"organization"`
	Dependencies map[string]Dependency  `yaml:"dependencies"`
	Support      map[string]Support     `yaml:"support"`
}

func GetNewProject(projectId, teamName, organizationName string) *Project {
	if teamName == "" {
		teamName = "Akatsuki"
	} 
	if organizationName == "" {
		organizationName = "Jailux"
	} 

	return &Project {
		ProjectId:    projectId,
		TeamName:     teamName,
		Organization: organizationName,
		Dependencies: nil,
		Support:      nil,
	}
}