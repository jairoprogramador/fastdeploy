package dto

type VariablesDTO []VariableDTO

type VariableDTO struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Value       string `yaml:"value"`
}
