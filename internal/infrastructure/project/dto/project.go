package dto

type ProjectDTO struct {
	ID           string `yaml:"id"`
	Name         string `yaml:"name"`
	Version      string `yaml:"version"`
	Team         string `yaml:"team"`
	Description  string `yaml:"description"`
	Organization string `yaml:"organization"`
}