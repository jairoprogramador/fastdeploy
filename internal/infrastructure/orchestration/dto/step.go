package dto

type StepRecordDTO struct {
	Name              string             `yaml:"name"`
	Status            string             `yaml:"status"`
	CommandExecutions []CommandRecordDTO `yaml:"commands"`
}
