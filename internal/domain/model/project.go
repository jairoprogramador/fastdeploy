package model

type Project struct {
	ProjectId    string                   `yaml:"projectId"`
	ProjectName  string                   `yaml:"projectName"`
	TeamName     string                   `yaml:"teamName"`
	Organization string                   `yaml:"organization"`
	Dependencies map[string]Dependency  `yaml:"dependencies"`
	Support      map[string]Support     `yaml:"support"`
}

func GetNewProject(projectId, projectName, teamName, organizationName string) *Project {
	if teamName == "" {
		teamName = "Akatsuki"
	} 
	if organizationName == "" {
		organizationName = "Jailux"
	} 

	return &Project {
		ProjectId:    projectId,
		ProjectName:  projectName,
		TeamName:     teamName,
		Organization: organizationName,
		Dependencies: nil,
		Support:      nil,
	}
}