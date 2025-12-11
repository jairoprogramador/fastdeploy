package vos

import "errors"

type StepName string

func NewStepName(value string) (StepName, error) {
	if value == "" {
		return "", errors.New("scope value cannot be empty")
	}
	return StepName(value), nil
}

func (e StepName) String() string {
	return string(e)
}
