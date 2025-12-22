package vos

import (
	"fmt"
	"strings"
)

const (
	varsExtension = "vars"
	stateExtension = "tb"
)

type FileName struct {
	value string
}

func NewFileName(name, extension string) (FileName, error) {
	if name == "" {
		return FileName{}, fmt.Errorf("base name for file cannot be empty")
	}
	if extension == "" {
		return FileName{}, fmt.Errorf("extension for file cannot be empty")
	}
	cleanExtension := strings.TrimPrefix(extension, ".")
	return FileName{value: fmt.Sprintf("%s.%s", name, cleanExtension)}, nil
}

func NewVarsFileName(scopeName string) (FileName, error) {
	return NewFileName(scopeName, varsExtension)
}

func NewStateFileName(scopeName string) (FileName, error) {
	return NewFileName(scopeName, stateExtension)
}

func (f FileName) String() string {
	return f.value
}
