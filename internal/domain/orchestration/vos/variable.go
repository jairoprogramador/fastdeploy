package vos

import "errors"

// Variable representa un par clave-valor en el mapa de variables compartido.
// Es un Objeto de Valor inmutable.
type Variable struct {
	key   string
	value string
}

// NewVariable crea un nuevo y validado Objeto de Valor Variable.
func NewVariable(key, value string) (Variable, error) {
	if key == "" {
		return Variable{}, errors.New("la clave de la variable no puede estar vacía")
	}
	if value == "" {
		return Variable{}, errors.New("el valor de la variable no puede estar vacío")
	}
	return Variable{
		key:   key,
		value: value,
	}, nil
}

// Key devuelve la clave de la variable.
func (v Variable) Key() string {
	return v.key
}

// Value devuelve el valor de la variable.
func (v Variable) Value() string {
	return v.value
}

// Equals comprueba la igualdad de valor con otra Variable.
func (v Variable) Equals(other Variable) bool {
	return v.key == other.key && v.value == other.value
}
