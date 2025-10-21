package vos

import (
	"errors"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"
)

type Output struct {
	name  string
	value string
}

func NewOutput(outputDef vos.Output) Output {
	return Output{
		name:  outputDef.Name(),
		value: outputDef.Probe(),
	}
}

func NewOutputFromNameAndValue(name, value string) (Output, error) {
	if name == "" {
		return Output{}, errors.New("el nombre de la salida no puede estar vacío")
	}
	if value == "" {
		return Output{}, errors.New("el valor de la salida no puede estar vacío")
	}

	return Output{
		name:  name,
		value: value,
	}, nil
}

func (v Output) Name() string {
	return v.name
}

func (v Output) Value() string {
	return v.value
}

func (v Output) Equals(other Output) bool {
	return v.name == other.name && v.value == other.value
}
