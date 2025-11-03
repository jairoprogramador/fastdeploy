package dto

type CommandDefinitionDTO struct {
	Name          string                `yaml:"name"`
	Description   string                `yaml:"description"`
	Cmd           string                `yaml:"cmd"`
	Workdir       string                `yaml:"workdir"`
	TemplateFiles []string              `yaml:"templates"`
	Outputs       []OutputDefinitionDTO `yaml:"outputs"`
}
