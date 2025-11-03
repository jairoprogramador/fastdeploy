package vos

import (
	"errors"
)

type CommandDefinition struct {
	name          string
	cmd           string
	workdir       string
	templateFiles []string
	outputs       []OutputDefinition
}

type CommandOption func(*CommandDefinition)

func NewCommandDefinition(name, cmd string, opts ...CommandOption) (CommandDefinition, error) {
	if name == "" {
		return CommandDefinition{}, errors.New("el nombre de la definición de comando no puede estar vacío")
	}
	if cmd == "" {
		return CommandDefinition{}, errors.New("el comando no puede estar vacío")
	}

	cmdDef := &CommandDefinition{
		name: name,
		cmd:  cmd,
	}

	for _, opt := range opts {
		opt(cmdDef)
	}

	if len(cmdDef.templateFiles) > 0 {
		templateFilesMap := make(map[string]struct{})
		for _, file := range cmdDef.templateFiles {
			if _, exists := templateFilesMap[file]; exists {
				return CommandDefinition{}, errors.New("archivo de plantilla duplicado")
			}
			templateFilesMap[file] = struct{}{}
		}
	}

	if len(cmdDef.outputs) > 0 {
		outputNames := make(map[string]struct{})
		for _, output := range cmdDef.outputs {
			if _, exists := outputNames[output.Name()]; exists {
				return CommandDefinition{}, errors.New("salida duplicada")
			}
			outputNames[output.Name()] = struct{}{}
		}
	}

	return *cmdDef, nil
}

func WithWorkdir(workdir string) CommandOption {
	return func(c *CommandDefinition) {
		c.workdir = workdir
	}
}

func WithTemplateFiles(files []string) CommandOption {
	return func(c *CommandDefinition) {
		c.templateFiles = files
	}
}

func WithOutputs(outputs []OutputDefinition) CommandOption {
	return func(c *CommandDefinition) {
		c.outputs = outputs
	}
}

func (cd CommandDefinition) Name() string {
	return cd.name
}

func (cd CommandDefinition) Cmd() string {
	return cd.cmd
}

func (cd CommandDefinition) Workdir() string {
	return cd.workdir
}

func (cd CommandDefinition) TemplateFiles() []string {
	filesCopy := make([]string, len(cd.templateFiles))
	copy(filesCopy, cd.templateFiles)
	return filesCopy
}

func (cd CommandDefinition) Outputs() []OutputDefinition {
	outputsCopy := make([]OutputDefinition, len(cd.outputs))
	copy(outputsCopy, cd.outputs)
	return outputsCopy
}
