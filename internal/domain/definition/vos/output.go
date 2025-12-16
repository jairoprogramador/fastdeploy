package vos

import "errors"

type OutputDefinition struct {
	name        string
	description string
	probe       string // Regex
}

func NewOutputDefinition(name, description, probe string) (OutputDefinition, error) {
	if name == "" {
		return OutputDefinition{}, errors.New("el nombre de la salida no puede estar vacío")
	}
	if probe == "" {
		return OutputDefinition{}, errors.New("la sonda (probe) de la salida no puede estar vacía")
	}
	return OutputDefinition{
		name:        name,
		description: description,
		probe:       probe,
	}, nil
}

func (o *OutputDefinition) Name() string {
	return o.name
}

func (o *OutputDefinition) Description() string {
	return o.description
}

func (o *OutputDefinition) Probe() string {
	return o.probe
}
