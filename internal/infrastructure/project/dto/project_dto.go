package dto

type ProjectDto struct {
	Organization string         `yaml:"organization"`
	Team         string         `yaml:"team"`
	Project      ProjectInfo    `yaml:"project"`
	Repository   RepositoryInfo `yaml:"repository"`
	Technology   TechnologyInfo `yaml:"technology"`
	Deployment   DeploymentInfo `yaml:"deployment"`
}
