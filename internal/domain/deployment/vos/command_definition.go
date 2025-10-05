package vos

import (
	"errors"
)

// CommandDefinition representa la definición estática de un comando a ejecutar.
// Es un Objeto de Valor inmutable.
type CommandDefinition struct {
	name          string
	description   string
	cmdTemplate   string
	workdir       string
	templateFiles []string
	outputs       []OutputProbe
}

// CommandOption define una firma para las opciones funcionales que configuran un CommandDefinition.
type CommandOption func(*CommandDefinition)

// NewCommandDefinition crea un nuevo y validado Objeto de Valor CommandDefinition.
// Utiliza el patrón de Opciones Funcionales para manejar los atributos opcionales.
func NewCommandDefinition(name, cmdTemplate string, opts ...CommandOption) (CommandDefinition, error) {
	if name == "" {
		return CommandDefinition{}, errors.New("el nombre de la definición de comando no puede estar vacío")
	}
	if cmdTemplate == "" {
		return CommandDefinition{}, errors.New("la plantilla de comando (cmdTemplate) no puede estar vacía")
	}

	cmdDef := &CommandDefinition{
		name:        name,
		cmdTemplate: cmdTemplate,
	}

	for _, opt := range opts {
		opt(cmdDef)
	}

	return *cmdDef, nil
}

// WithDescription establece la descripción para un CommandDefinition.
func WithDescription(description string) CommandOption {
	return func(c *CommandDefinition) {
		c.description = description
	}
}

// WithWorkdir establece el directorio de trabajo para un CommandDefinition.
func WithWorkdir(workdir string) CommandOption {
	return func(c *CommandDefinition) {
		c.workdir = workdir
	}
}

// WithTemplateFiles establece los archivos de plantilla para un CommandDefinition.
func WithTemplateFiles(files []string) CommandOption {
	return func(c *CommandDefinition) {
		c.templateFiles = files
	}
}

// WithOutputs establece las sondas de salida para un CommandDefinition.
func WithOutputs(outputs []OutputProbe) CommandOption {
	return func(c *CommandDefinition) {
		c.outputs = outputs
	}
}

// Name devuelve el nombre del comando.
func (cd CommandDefinition) Name() string {
	return cd.name
}

// Description devuelve la descripción del comando.
func (cd CommandDefinition) Description() string {
	return cd.description
}

// CmdTemplate devuelve la plantilla de comando.
func (cd CommandDefinition) CmdTemplate() string {
	return cd.cmdTemplate
}

// Workdir devuelve el directorio de trabajo del comando.
func (cd CommandDefinition) Workdir() string {
	return cd.workdir
}

// TemplateFiles devuelve la lista de archivos de plantilla asociados al comando.
func (cd CommandDefinition) TemplateFiles() []string {
	// Devolvemos una copia para proteger la inmutabilidad.
	filesCopy := make([]string, len(cd.templateFiles))
	copy(filesCopy, cd.templateFiles)
	return filesCopy
}

// Outputs devuelve la lista de sondas de salida para el comando.
func (cd CommandDefinition) Outputs() []OutputProbe {
	// Devolvemos una copia para proteger la inmutabilidad.
	outputsCopy := make([]OutputProbe, len(cd.outputs))
	copy(outputsCopy, cd.outputs)
	return outputsCopy
}
