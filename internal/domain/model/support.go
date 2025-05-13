package model

type Support struct {
	Type    string            `yaml:"type"`
	Name    string            `yaml:"name"`
	Version string            `yaml:"version"`
	URL     string            `yaml:"url"`
	Config  map[string]string `yaml:"config"`
}
