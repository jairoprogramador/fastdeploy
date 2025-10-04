package dto

// DTOs para el unmarshalling de YAML.
// Desacoplan el modelo de dominio de la estructura de los archivos de configuración.
type CommandDefinitionDTO struct {
	Name          string           `yaml:"name"`
	Description   string           `yaml:"description"`
	Cmd           string           `yaml:"cmd"`
	Workdir       string           `yaml:"workdir"`
	TemplateFiles []string         `yaml:"templates"`
	Outputs       []OutputProbeDTO `yaml:"outputs"`
}