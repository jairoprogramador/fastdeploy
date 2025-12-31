package dto

type FdConfigDTO struct {
	Project  ProjectDTO  `yaml:"project"`
	Template TemplateDTO `yaml:"template"`
}