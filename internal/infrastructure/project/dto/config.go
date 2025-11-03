package dto

type ConfigDTO struct {
	Project  ProjectDTO  `yaml:"project"`
	Template TemplateDTO `yaml:"template"`
	State    struct {
		Backend string `yaml:"backend"`
		URL     string `yaml:"url"`
	} `yaml:"state"`
}
