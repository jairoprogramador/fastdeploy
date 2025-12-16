package vos

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var stepNameRegex = regexp.MustCompile(`^(\d+)-(.+)$`)

type StepNameDefinition struct {
	order int
	name  string
}

func NewStepNameDefinition(dirName string) (StepNameDefinition, error) {
	matches := stepNameRegex.FindStringSubmatch(dirName)
	if len(matches) != 3 {
		return StepNameDefinition{}, fmt.Errorf("el nombre del directorio del paso '%s' no sigue el formato 'NN-nombre'", dirName)
	}

	order, err := strconv.Atoi(matches[1])
	if err != nil {
		return StepNameDefinition{}, fmt.Errorf("no se pudo parsear el número de orden del paso '%s'", dirName)
	}

	name := matches[2]
	if name == "" {
		return StepNameDefinition{}, errors.New("el nombre del paso no puede estar vacío")
	}

	return StepNameDefinition{order: order, name: name}, nil
}

func (s StepNameDefinition) Order() int {
	return s.order
}

func (s StepNameDefinition) Name() string {
	return s.name
}

func (s StepNameDefinition) FullName() string {
	return fmt.Sprintf("%02d-%s", s.order, s.name)
}
