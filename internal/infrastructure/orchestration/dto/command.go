package dto

type CommandRecordDTO struct {
	Name        string            `yaml:"name"`
	Status      string            `yaml:"status"`
	ResolvedCmd string            `yaml:"cmd"`
	Record      string            `yaml:"record"`
	OutputVars  map[string]string `yaml:"outputs"`
}
