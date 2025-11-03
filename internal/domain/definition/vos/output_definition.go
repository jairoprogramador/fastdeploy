package vos

import (
	"errors"
	"fmt"
	"regexp"
)

type OutputDefinition struct {
	name  string
	probe string
}

func NewOutputDefinition(name, probe string) (OutputDefinition, error) {
	if probe == "" {
		return OutputDefinition{}, errors.New("la expresión de la sonda de salida no puede estar vacía")
	}
	if _, err := regexp.Compile(probe); err != nil {
		return OutputDefinition{}, fmt.Errorf("la expresión regular de la sonda no es válida: %w", err)
	}

	return OutputDefinition{
		name:  name,
		probe: probe,
	}, nil
}

func (op OutputDefinition) Name() string {
	return op.name
}

func (op OutputDefinition) Probe() string {
	return op.probe
}
