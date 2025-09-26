package values

import (
	"errors"
	"strings"
	"github.com/jairoprogramador/fastdeploy/internal/domain/variable/values"
)

const ERROR_COMMAND_CANNOT_BE_EMPTY = "command cannot be empty"
const DEFAULT_WORKDIR = "."

var DEFAULT_TEMPLATES = []string{}
var DEFAULT_OUTPUTS = []OutputValue{}

type CommandValue struct {
	name    string
	command string
	workdir string
	outputs []OutputValue
	/* timeout int
	outputDir string
	encoding string
	shell string
	async bool
	logLevel string
	onFailure string
	inputStream string */

	templates []string
}

/* func DefaultCommand(name, command string) (CommandValue, error) {
	name = strings.TrimSpace(name)
	command = strings.TrimSpace(command)

	if command == "" {
		return CommandValue{}, errors.New(ERROR_COMMAND_CANNOT_BE_EMPTY)
	}

	if name == "" {
		name = "command name not specified"
	}

	return CommandValue{
		name:      name,
		command:   command,
		outputs:   DEFAULT_OUTPUTS,
		templates: DEFAULT_TEMPLATES,
		workdir:   DEFAULT_WORKDIR,
	}, nil
} */

func NewCommand(
	name, command, workdir string,
	outputs []OutputValue,
	templating []string) (CommandValue, error) {

	name = strings.TrimSpace(name)
	command = strings.TrimSpace(command)
	workdir = strings.TrimSpace(workdir)

	if command == "" {
		return CommandValue{}, errors.New(ERROR_COMMAND_CANNOT_BE_EMPTY)
	}

	if templating == nil {
		templating = DEFAULT_TEMPLATES
	}

	if outputs == nil {
		outputs = DEFAULT_OUTPUTS
	}

	if workdir == "" {
		workdir = DEFAULT_WORKDIR
	}

	return CommandValue{
		name:      name,
		command:   command,
		outputs:   outputs,
		templates: templating,
		workdir:   workdir,
	}, nil
}

func (c CommandValue) GetCommand() string {
	return c.command
}

func (c CommandValue) GetWorkdir() string {
	return c.workdir
}

func (c CommandValue) GetTemplates() []string {
	return c.templates
}

func (c CommandValue) GetName() string {
	return c.name
}

func (c *CommandValue) AddOutput(output OutputValue) {
	c.outputs = append(c.outputs, output)
}

func (c CommandValue) IsValid(outputCommand string, errCommand error) ([]values.VariableValue, error) {
	if errCommand != nil {
		return []values.VariableValue{}, errCommand
	}

	var variables = make([]values.VariableValue, 0)

	for _, output := range c.outputs {
		variable, err := output.IsValid(outputCommand)
		if err != nil {
			return []values.VariableValue{}, err
		}
		if variable.GetValue() != "" && variable.GetName() != "" {
			variables = append(variables, variable)
		}
	}

	return variables, nil
}
/*
func (c CommandValue) Execute(
	templateCommand port.TemplateCommand,
	executorCommand port.ExecutorCommand) ([]VariableValue, error) {

	if err := templateCommand.Process(c.templates); err != nil {
		return []VariableValue{}, err
	}

	variables, err := c.IsValid(executorCommand.Run(c.command, c.workdir))
	if err != nil {
		return []VariableValue{}, err
	}

	return variables, nil
} */
