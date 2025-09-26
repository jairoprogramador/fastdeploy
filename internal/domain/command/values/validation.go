package values

import (
	"errors"
	"regexp"
	"strings"
)

var ansiEscapeCodeRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
var defaultAllMatches = [][]string{}

type validationValueType int

const (
	DEFAULT  validationValueType = iota // Caso 1: fallo por exit code != 0
	REGEX                                 // Caso 2: fallo si regex no coincide
	CONTAINS                              // Caso 2: fallo si no contiene texto
	CUSTOM                                // Para futuras validaciones personalizadas
)

type ValidationValue struct {
	Type        validationValueType
	probe       string
	InvertMatch bool
}

func NewValidation(probe string, typeValidation int) *ValidationValue {
	probe = strings.TrimSpace(probe)

	if probe == "" {
		return DefaultValidation()
	}

	if typeValidation > 3 {
		return DefaultValidation()
	}

	return &ValidationValue{
		Type:  validationValueType(typeValidation),
		probe: probe,
	}
}

func DefaultValidation() *ValidationValue {
	return &ValidationValue{
		Type: DEFAULT,
	}
}

func RegexValidation(pattern string) *ValidationValue {
	pattern = strings.TrimSpace(pattern)

	if pattern == "" {
		return DefaultValidation()
	}

	return &ValidationValue{
		Type:  REGEX,
		probe: pattern,
	}
}

func ContainsValidation(expected string) *ValidationValue {
	expected = strings.TrimSpace(expected)

	if expected == "" {
		return DefaultValidation()
	}

	return &ValidationValue{
		Type:  CONTAINS,
		probe: expected,
	}
}

func (e *ValidationValue) IsValid(outputCommand string) ([][]string, error) {
	if e.Type == REGEX {
		allMatches, err := e.GetAllSubmatch(outputCommand)
		if err != nil {
			return defaultAllMatches, err
		}

		if len(allMatches) == 0 {
			return defaultAllMatches, errors.New("no matches found for regex: " + e.probe)
		}

		return allMatches, nil
	}

	return defaultAllMatches, nil
}

func (e *ValidationValue) GetAllSubmatch(outputCommand string) ([][]string, error) {
	regexpValidation, err := regexp.Compile(e.probe)
	if err != nil {
		return defaultAllMatches, err
	}

	outputCommandNoEscapeCode := ansiEscapeCodeRegex.ReplaceAllString(outputCommand, "")

	allSubmatch := regexpValidation.FindAllStringSubmatch(outputCommandNoEscapeCode, -1)

	return allSubmatch, nil
}
