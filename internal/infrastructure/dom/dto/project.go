package dto

type ProjectDTO struct {
	ID          string `yaml:"project_id"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Team        string `yaml:"team"`
}
