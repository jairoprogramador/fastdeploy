package vos

import "errors"

type VariableDefinition struct {
	name  string
	value string
}

func NewVariableDefinition(name, value string) (VariableDefinition, error) {
	if name == "" {
		return VariableDefinition{}, errors.New("el nombre de la variable no puede estar vacío")
	}

	if value == "" {
		return VariableDefinition{}, errors.New("el valor de la variable no puede estar vacío")
	}

	return VariableDefinition{
		name:  name,
		value: value,
	}, nil
}

func (v VariableDefinition) Name() string {
	return v.name
}

func (v VariableDefinition) Value() string {
	return v.value
}
