package entities

import (
	"errors"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

type Step struct {
	workspaceRoot string
	name          string
	commands      []vos.Command
	variables     []vos.VariableSet
}

type StepOption func(*Step)

func NewStep(name string, opts ...StepOption) (Step, error) {
	if name == "" {
		return Step{}, errors.New("el nombre del paso no puede estar vacío")
	}

	step := &Step{
		name: name,
	}

	for _, opt := range opts {
		opt(step)
	}

	if len(step.commands) == 0 {
		return Step{}, errors.New("un paso debe tener al menos un comando")
	}

	return *step, nil
}

func WithCommands(commands []vos.Command) StepOption {
	return func(s *Step) {
		s.commands = commands
	}
}

func WithVariables(variables []vos.VariableSet) StepOption {
	return func(s *Step) {
		s.variables = variables
	}
}

func WithWorkspaceRoot(workspaceRoot string) StepOption {
	return func(s *Step) {
		s.workspaceRoot = workspaceRoot
	}
}

func (sd Step) Name() string {
	return sd.name
}

func (sd Step) WorkspaceRoot() string {
	return sd.workspaceRoot
}

func (sd Step) Commands() []vos.Command {
	commandsCopy := make([]vos.Command, len(sd.commands))
	copy(commandsCopy, sd.commands)
	return commandsCopy
}

func (sd Step) Variables() []vos.VariableSet {
	variablesCopy := make([]vos.VariableSet, len(sd.variables))
	copy(variablesCopy, sd.variables)
	return variablesCopy
}
