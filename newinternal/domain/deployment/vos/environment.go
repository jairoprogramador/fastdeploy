package vos

import (
	"errors"
)

// Environment representa un ambiente de despliegue.
// Es un Objeto de Valor, identificado por sus atributos, no por un ID específico.
// Su estado es inmutable después de la creación.
type Environment struct {
	name        string
	description string
	value       string
}

// NewEnvironment crea un nuevo y validado Objeto de Valor Environment.
// El constructor asegura que el objeto siempre se encuentre en un estado válido.
func NewEnvironment(name, description, value string) (Environment, error) {
	if name == "" {
		return Environment{}, errors.New("el nombre del ambiente no puede estar vacío")
	}
	if value == "" {
		return Environment{}, errors.New("el valor del ambiente no puede estar vacío")
	}

	return Environment{
		name:        name,
		description: description,
		value:       value,
	}, nil
}

// Name devuelve el nombre del ambiente.
func (e Environment) Name() string {
	return e.name
}

// Description devuelve la descripción del ambiente.
func (e Environment) Description() string {
	return e.description
}

// Value devuelve el identificador corto del ambiente.
func (e Environment) Value() string {
	return e.value
}

// Equals comprueba la igualdad de valor con otro objeto Environment.
func (e Environment) Equals(other Environment) bool {
	return e.name == other.name && e.value == other.value
}