package dto

// DTOs para el unmarshalling de YAML.
// Desacoplan el modelo de dominio de la estructura de los archivos de configuración.
type EnvironmentDTO struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}
