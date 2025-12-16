package vos

import "errors"

type VariableDefinition struct {
	name  string
	value interface{}
}

func NewVariableDefinition(name string, value interface{}) (VariableDefinition, error) {
	if name == "" {
		return VariableDefinition{}, errors.New("el nombre de la variable no puede estar vacío")
	}
	return VariableDefinition{name: name, value: value}, nil
}

func (v VariableDefinition) Name() string {
	return v.name
}

func (v VariableDefinition) Value() interface{} {
	return v.value
}
