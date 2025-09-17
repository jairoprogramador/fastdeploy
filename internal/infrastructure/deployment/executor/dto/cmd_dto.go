package dto

type CmdDTO struct {
	Name            string               `yaml:"name"`
	Description     string               `yaml:"description,omitempty"`
	Cmd             string               `yaml:"cmd"`
	ContinueOnError bool                 `yaml:"continue_on_error,omitempty"`
	Workdir         string               `yaml:"workdir,omitempty"`
	Result          string               `yaml:"result,omitempty"`
	Variables       []VariablePatternDTO `yaml:"variables,omitempty"`
}
