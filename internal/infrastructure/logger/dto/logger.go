package dto

import "time"

type LoggerDTO struct {
	Status    string            `yaml:"status"`
	StartTime time.Time         `yaml:"start_time"`
	EndTime   time.Time         `yaml:"end_time,omitempty"`
	Steps     []StepDTO         `yaml:"steps"`
	Context   map[string]string `yaml:"context"`
	Revision  string            `yaml:"revision"`
}
