package project

type ProjectEntity struct {
	Organization string `yaml:"organization"`
	ProjectID    string `yaml:"projectId"`
	ProjectName  string `yaml:"projectName"`
	Repository   string `yaml:"repository"`
	Technology   string `yaml:"technology"`
	Version      string `yaml:"version"`
	TeamName     string `yaml:"teamName"`
}
