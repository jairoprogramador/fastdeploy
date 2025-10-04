package dto

// DTOs para el unmarshalling de YAML.
// Desacoplan el modelo de dominio de la estructura de los archivos de configuración.
type OutputProbeDTO struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Probe       string `yaml:"probe"`
}