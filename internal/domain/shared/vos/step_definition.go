package vos

type StepDefinition struct {
	Name      string
	Cmd       string
	Workdir   string
	Templates []*FileTemplate
	Outputs   []*Output
}

func (sd *StepDefinition) IsTest() bool {
	return sd.Name == "test"
}
