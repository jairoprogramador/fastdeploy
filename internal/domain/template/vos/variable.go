package vos

import "errors"

type Variable struct {
	name string
	value string
}

func NewVariable(name, value string) (Variable, error) {
	if name == "" {
		return Variable{}, errors.New("el nombre de la variable no puede estar vacío")
	}

	if value == "" {
		return Variable{}, errors.New("el valor de la variable no puede estar vacío")
	}

	return Variable{
		name: name,
		value: value,
	}, nil
}

func (v Variable) Name() string {
	return v.name
}

func (v Variable) Value() string {
	return v.value
}