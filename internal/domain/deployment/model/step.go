package model

type TypeStep string

const (
	Command   TypeStep = "command"
	Container TypeStep = "container"
	Check     TypeStep = "check"
)

type Step struct {
	Name      string     `yaml:"name"`
	Type      string     `yaml:"type"`
	Command   string     `yaml:"command,omitempty"`
	Timeout   string     `yaml:"timeout,omitempty"`
	Retry     *Retry     `yaml:"retry,omitempty"`
	If        string     `yaml:"if,omitempty"`
	Then      string     `yaml:"then,omitempty"`
	Skip      []string   `yaml:"skip,omitempty"`
	Parallel  []Step     `yaml:"parallel,omitempty"`
	Variables []Variable `yaml:"variables,omitempty"`
}
