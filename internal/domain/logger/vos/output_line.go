package vos

import "time"

type OutputLine struct {
	timestamp time.Time `yaml:"timestamp"`
	line      string    `yaml:"line"`
}

func NewOutputLine(line string) OutputLine {
	return OutputLine{
		timestamp: time.Now(),
		line:      line,
	}
}

func (o OutputLine) Timestamp() time.Time {
	return o.timestamp
}

func (o OutputLine) Line() string {
	return o.line
}
