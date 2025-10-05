package vos

import "errors"

type StateStep struct {
	name string
	successful bool
}

func NewStateStep(name string, successful bool) (StateStep, error) {
	if name == "" {
		return StateStep{}, errors.New("name cannot be empty")
	}
	return StateStep{name: name, successful: successful}, nil
}

func (s StateStep) GetName() string {
	return s.name
}

func (s StateStep) IsSuccessful() bool {
	return s.successful
}

func (s StateStep) Equals(other StateStep) bool {
	return s.name == other.name && s.successful == other.successful
}