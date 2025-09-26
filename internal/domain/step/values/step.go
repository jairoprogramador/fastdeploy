package values

import (
	"errors"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/command/values"
)

type StepValue struct {
	name     string
	commands []values.CommandValue
}

func NewStepValue(name string, commands []values.CommandValue) (StepValue, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return StepValue{}, errors.New("name cannot be empty")
	}

	if len(commands) == 0 {
		return StepValue{}, errors.New("commands cannot be empty")
	}

	return StepValue{name: name, commands: commands}, nil
}

func (s StepValue) GetName() string {
	return s.name
}

func (s StepValue) GetCommands() []values.CommandValue {
	return s.commands
}
