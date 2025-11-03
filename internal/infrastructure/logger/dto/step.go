package dto

import "time"

type StepDTO struct {
	Name string `yaml:"name"`
	Status string `yaml:"status"`
	StartTime time.Time `yaml:"start_time"`
	EndTime time.Time `yaml:"end_time,omitempty"`
	Reason string `yaml:"reason,omitempty"`
	Tasks []TaskDTO `yaml:"tasks"`
	Err error `yaml:"err,omitempty"`
}