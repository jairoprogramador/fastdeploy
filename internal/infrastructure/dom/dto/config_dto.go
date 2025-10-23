package dto

type FileConfig struct {
	Project ProjectDTO `yaml:"project"`
	Template TemplateDTO `yaml:"template"`
	Technology TechnologyDTO `yaml:"technology"`
	Runtime RuntimeDTO `yaml:"runtime"`
	State struct {
		Backend string `yaml:"backend"`
		URL     string `yaml:"url"`
	} `yaml:"state"`
}
