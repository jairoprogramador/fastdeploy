package model

type DeploymentEntity struct {
	Version     string    `yaml:"version"`
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Variables   Variables `yaml:"variables"`
	Steps       []Step    `yaml:"steps"`
}

func (d *DeploymentEntity) HasType(typeStep string) bool {
	for _, step := range d.Steps {
		if step.Type == typeStep {
			return true
		}
	}
	return false
}
