package model

type Retry struct {
	Attempts int    `yaml:"attempts"`
	Delay    string `yaml:"delay"`
}
