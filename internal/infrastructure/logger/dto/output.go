package dto

import "time"

type OutputDTO struct {
	Timestamp time.Time `yaml:"timestamp"`
	Line string `yaml:"line"`
}