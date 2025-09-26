package values

import (
	"errors"
	"strings"
	"github.com/jairoprogramador/fastdeploy/internal/domain/variable/values"
)

type OutputValue struct {
	name        string
	description string
	validator   *ValidationValue
}

func NewOutput(name, description string, validator *ValidationValue) (OutputValue, error) {

	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if validator == nil {
		return OutputValue{}, errors.New("validation cannot be empty")
	}

	return OutputValue{name: name, description: description, validator: validator}, nil
}

func (o OutputValue) NameIsEmpty() bool {
	return o.name == ""
}

func (o OutputValue) GetName() string {
	return o.name
}

func (o OutputValue) GetDescription() string {
	return o.description
}

func (o OutputValue) IsValid(outputCommand string) (values.VariableValue, error) {
	allMatches, err := o.validator.IsValid(outputCommand)
	if err != nil {
		return values.VariableValue{}, err
	}

	value := o.getFirstFoundValue(allMatches)

	return values.NewVariable(o.name, value), nil
}

func (e *OutputValue) getFirstFoundValue(allMatches [][]string) string {
	if len(allMatches) == 0 {
		return ""
	}

	return allMatches[0][1]
}
