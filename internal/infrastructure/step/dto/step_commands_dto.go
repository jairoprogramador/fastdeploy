package dto

type StepCommandsDTO []StepCommandDTO

type StepCommandDTO struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Cmd         string `yaml:"cmd"`
	//ContinueOnError bool              `yaml:"continue_on_error,omitempty"`
	Workdir string `yaml:"workdir,omitempty"`
	//Result          string            `yaml:"result,omitempty"`
	Outputs    []StepOutputDTO  `yaml:"outputs,omitempty"`
	Templates StepTemplatesDTO `yaml:"templates,omitempty"`
	//NotExecuteLocal bool              `yaml:"not_execute_local,omitempty"`
}
