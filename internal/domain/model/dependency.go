package model

type Dependency struct {
	Type     string            `yaml:"type"`
	Name     string            `yaml:"name"`
	Version  string            `yaml:"version"`
	Required bool              `yaml:"required"`
	Config   map[string]string `yaml:"config"`
}

func GetNewDependency(typeDependency, name, version string, required bool, config map[string]string) Dependency {
	return Dependency {
		Type: typeDependency,
		Name: name,
		Version: version,
		Required: required,
		Config: config,
	}
}