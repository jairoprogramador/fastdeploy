package model

type Deployment struct {
	Version     string    `yaml:"version"`
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Variables   Variables `yaml:"variables"`
	Steps       []Step    `yaml:"steps"`
}

type Variables struct {
	Global []Variable `yaml:"global"`
	Local  []Variable `yaml:"local"`
}

type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Step struct {
	Name      string       `yaml:"name"`
	Type      string       `yaml:"type"`
	Command   string       `yaml:"command,omitempty"`
	Timeout   string       `yaml:"timeout,omitempty"`
	Retry     *RetryConfig `yaml:"retry,omitempty"`
	If        string       `yaml:"if,omitempty"`
	Then      string       `yaml:"then,omitempty"`
	Skip      []string     `yaml:"skip,omitempty"`
	Parallel  []Step       `yaml:"parallel,omitempty"`
	Variables []Variable   `yaml:"variables,omitempty"`
}

type RetryConfig struct {
	Attempts int    `yaml:"attempts"`
	Delay    string `yaml:"delay"`
}
