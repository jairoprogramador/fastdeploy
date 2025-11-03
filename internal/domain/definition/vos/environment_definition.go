package vos

import (
	"errors"
)

type EnvironmentDefinition struct {
	name  string
	value string
}

func NewEnvironmentDefinition(name, value string) (EnvironmentDefinition, error) {
	if name == "" {
		return EnvironmentDefinition{}, errors.New("el nombre del ambiente no puede estar vacío")
	}
	if value == "" {
		return EnvironmentDefinition{}, errors.New("el valor del ambiente no puede estar vacío")
	}

	return EnvironmentDefinition{
		name:  name,
		value: value,
	}, nil
}

func (e EnvironmentDefinition) Name() string {
	return e.name
}

func (e EnvironmentDefinition) Value() string {
	return e.value
}
