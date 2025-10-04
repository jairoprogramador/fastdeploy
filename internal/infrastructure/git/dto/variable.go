package dto

// VariableDTO representa la estructura de una variable en un archivo YAML.
type VariableDTO struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
