package dto

type OutputDto struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Regex       string `yaml:"regex,omitempty"`
}