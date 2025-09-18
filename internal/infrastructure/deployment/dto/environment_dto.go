package dto

type EnvironmentDto struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}
