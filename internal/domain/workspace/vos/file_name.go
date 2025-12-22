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

func newFileName(name, extension string) (FileName, error) {
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
	return newFileName(scopeName, varsExtension)
}

func NewStateFileName(scopeName string) (FileName, error) {
	return newFileName(scopeName, stateExtension)
}

func (f FileName) String() string {
	return f.value
}
