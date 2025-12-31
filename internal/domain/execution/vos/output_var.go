package vos

import "errors"

const (
	// SharedScope es el identificador para las variables que son compartidas
	// entre todos los entornos.
	SharedScope = "shared"
)

type OutputVar struct {
	name     string
	value    string
	isShared bool
}

func NewOutputVar(name, value string, isShared bool) (OutputVar, error) {
	if name == "" {
		return OutputVar{}, errors.New("el nombre de la variable generada no puede estar vacío")
	}
	if value == "" {
		return OutputVar{}, errors.New("el valor de la variable generada no puede estar vacío")
	}

	return OutputVar{
		isShared: isShared,
		name:     name,
		value:    value,
	}, nil
}

func (ve *OutputVar) IsShared() bool {
	return ve.isShared
}

func (ve *OutputVar) Name() string {
	return ve.name
}

func (ve *OutputVar) Value() string {
	return ve.value
}
