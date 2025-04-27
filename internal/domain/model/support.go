package model

type Support struct {
	Type    string            `yaml:"type"`
	Name    string            `yaml:"name"`
	Version string            `yaml:"version"`
	URL     string            `yaml:"url"`
	Config  map[string]string `yaml:"config"`
}

func GetNewSupport(typeSupport, name, version, url string, config map[string]string) Support {
	return Support {
		Type: typeSupport,
		Name: name,
		Version: version,
		URL: url,
		Config: config,
	}
}