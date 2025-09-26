package dto

type StepOutputDTO struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Probe       string `yaml:"probe,omitempty"`
}
