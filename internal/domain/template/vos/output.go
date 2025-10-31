package vos

import (
	"errors"
	"fmt"
	"regexp"
)

type Output struct {
	name  string
	probe string
}

func NewOutput(name, probe string) (Output, error) {
	if probe == "" {
		return Output{}, errors.New("la expresión de la sonda de salida no puede estar vacía")
	}
	if _, err := regexp.Compile(probe); err != nil {
		return Output{}, fmt.Errorf("la expresión regular de la sonda no es válida: %w", err)
	}

	return Output{
		name:  name,
		probe: probe,
	}, nil
}

func (op Output) Name() string {
	return op.name
}

func (op Output) Probe() string {
	return op.probe
}
