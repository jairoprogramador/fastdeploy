package dto

// VariableDTO define una variable con su nombre, descripción y valor.
type VariableDTO struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Value       string `yaml:"value"`
}
