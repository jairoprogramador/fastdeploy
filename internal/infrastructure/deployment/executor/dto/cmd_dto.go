package dto

type CmdDto struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Cmd         string `yaml:"cmd"`
	ContinueOnError bool `yaml:"continue_on_error,omitempty"`
	Workdir     string `yaml:"workdir,omitempty"`
	Outputs     []OutputDto `yaml:"outputs,omitempty"`
}
