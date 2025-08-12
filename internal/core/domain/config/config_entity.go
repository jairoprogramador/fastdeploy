package config

type ConfigEntity struct {
	Organization string `yaml:"organization"`
	TeamName     string `yaml:"teamName"`
	Repository   string `yaml:"repository"`
}
