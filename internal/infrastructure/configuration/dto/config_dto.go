package dto

type ConfigDto struct {
	Organization string         `yaml:"organization"`
	Team         string         `yaml:"team"`
	Repository   string         `yaml:"repository"`
	Technology   TechnologyInfo `yaml:"technology"`
}
