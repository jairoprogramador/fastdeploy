package dto

type CmdDto struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Cmd         string `yaml:"cmd"`
	Dir         string `yaml:"dir,omitempty"`
	Outputs     []OutputDto `yaml:"outputs,omitempty"`
}
