package dto

type DomDTO struct {
	Product    ProductDTO    `yaml:"product"`
	Project    ProjectDTO    `yaml:"project"`
	Template   TemplateDTO   `yaml:"template"`
	Technology TechnologyDTO `yaml:"technology"`
}
