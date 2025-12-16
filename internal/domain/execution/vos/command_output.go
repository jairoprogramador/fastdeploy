package vos

import (
	"errors"
)

type CommandOutput struct {
	name  string
	probe string
}

func NewCommandOutput(name, probe string) (CommandOutput, error) {
	if name == "" {
		return CommandOutput{}, errors.New("el nombre de la sonda de comando no puede estar vacío")
	}
	if probe == "" {
		return CommandOutput{}, errors.New("la expresión de la sonda de comando no puede estar vacía")
	}

	return CommandOutput{
		name:  name,
		probe: probe,
	}, nil
}

func (op CommandOutput) Name() string {
	return op.name
}

func (op CommandOutput) Probe() string {
	return op.probe
}
