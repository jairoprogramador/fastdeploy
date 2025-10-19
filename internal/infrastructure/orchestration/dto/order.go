package dto

type OrderDTO struct {
	ID                string            `yaml:"id"`
	Status            string            `yaml:"status"`
	TargetEnvironment string            `yaml:"environment"`
	StepRecords       []StepRecordDTO   `yaml:"steps"`
	VariableMap       map[string]string `yaml:"variables"`
}
