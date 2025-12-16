package vos

import "errors"

type EnvironmentDefinition struct {
	name  string
	value string
}

func NewEnvironment(value, name string) (EnvironmentDefinition, error) {
	if value == "" {
		return EnvironmentDefinition{}, errors.New("el nombre del entorno no puede estar vacío")
	}
	if name == "" {
		return EnvironmentDefinition{}, errors.New("el nombre del entorno no puede estar vacío")
	}
	return EnvironmentDefinition{value: value}, nil
}

func (e EnvironmentDefinition) String() string {
	return e.value
}

func (e EnvironmentDefinition) Name() string {
	return e.name
}

func (e EnvironmentDefinition) Equals(other EnvironmentDefinition) bool {
	return e.value == other.value
}
