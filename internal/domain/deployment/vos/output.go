package vos

import (
	"errors"
	"fmt"
	"regexp"
)

// OutputProbe define una expresión regular para analizar la salida de la ejecución
// de un comando y extraer nuevas variables. Es un Objeto de Valor.
type OutputProbe struct {
	name        string // Nombre de la variable a crear. Puede estar vacío.
	description string
	probe       string // La expresión regular.
}

// NewOutputProbe crea un nuevo y validado OutputProbe.
func NewOutputProbe(name, description, probe string) (OutputProbe, error) {
	if probe == "" {
		return OutputProbe{}, errors.New("la expresión de la sonda de salida no puede estar vacía")
	}
	if _, err := regexp.Compile(probe); err != nil {
		return OutputProbe{}, fmt.Errorf("la expresión regular de la sonda no es válida: %w", err)
	}

	return OutputProbe{
		name:        name,
		description: description,
		probe:       probe,
	}, nil
}

// Name devuelve el nombre de la variable a crear.
func (op OutputProbe) Name() string {
	return op.name
}

// Description devuelve la descripción de la sonda.
func (op OutputProbe) Description() string {
	return op.description
}

// Probe devuelve la expresión regular de la sonda.
func (op OutputProbe) Probe() string {
	return op.probe
}
