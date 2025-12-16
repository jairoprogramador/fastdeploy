package vos

import (
	"errors"
)

type Command struct {
	name          string
	cmd           string
	workdir       string
	templateFiles []string
	outputs       []CommandOutput
}

type CommandOption func(*Command)

func NewCommand(name, cmd string, opts ...CommandOption) (Command, error) {
	if name == "" {
		return Command{}, errors.New("el nombre de comando no puede estar vacío")
	}
	if cmd == "" {
		return Command{}, errors.New("el comando no puede estar vacío")
	}

	cmdDef := &Command{
		name: name,
		cmd:  cmd,
	}

	for _, opt := range opts {
		opt(cmdDef)
	}

	return *cmdDef, nil
}

func WithWorkdir(workdir string) CommandOption {
	return func(c *Command) {
		c.workdir = workdir
	}
}

func WithTemplateFiles(files []string) CommandOption {
	return func(c *Command) {
		c.templateFiles = files
	}
}

func WithOutputs(outputs []CommandOutput) CommandOption {
	return func(c *Command) {
		c.outputs = outputs
	}
}

func (cd Command) Name() string {
	return cd.name
}

func (cd Command) Cmd() string {
	return cd.cmd
}

func (cd Command) Workdir() string {
	return cd.workdir
}

func (cd Command) TemplateFiles() []string {
	filesCopy := make([]string, len(cd.templateFiles))
	copy(filesCopy, cd.templateFiles)
	return filesCopy
}

func (cd Command) Outputs() []CommandOutput {
	outputsCopy := make([]CommandOutput, len(cd.outputs))
	copy(outputsCopy, cd.outputs)
	return outputsCopy
}
