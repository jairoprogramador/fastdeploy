package dto

import "time"

type TaskDTO struct {
	Name string `yaml:"name"`
	Status string `yaml:"status"`
	Command string `yaml:"command"`
	StartTime time.Time `yaml:"start_time"`
	EndTime time.Time `yaml:"end_time,omitempty"`
	Output []OutputDTO `yaml:"output,omitempty"`
	Err error `yaml:"err,omitempty"`
}