package config

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	Organization string `yaml:"organization"`
	ProjectID    string `yaml:"projectId"`
	ProjectName  string `yaml:"projectName"`
	Repository   string `yaml:"repository"`
	Technology   string `yaml:"technology"`
	Version      string `yaml:"version"`
	TeamName     string `yaml:"teamName"`
}

func (c *Config) ToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
