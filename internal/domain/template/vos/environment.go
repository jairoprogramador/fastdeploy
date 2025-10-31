package vos

import (
	"errors"
)

type Environment struct {
	name        string
	value       string
}

func NewEnvironment(name, value string) (Environment, error) {
	if name == "" {
		return Environment{}, errors.New("el nombre del ambiente no puede estar vacío")
	}
	if value == "" {
		return Environment{}, errors.New("el valor del ambiente no puede estar vacío")
	}

	return Environment{
		name:        name,
		value:       value,
	}, nil
}

func RehydrateEnvironment(value string) Environment {
	return Environment{
		name:        value,
		value:       value,
	}
}

func (e Environment) Name() string {
	return e.name
}

func (e Environment) Value() string {
	return e.value
}

func (e Environment) Equals(other Environment) bool {
	return e.name == other.name && e.value == other.value
}