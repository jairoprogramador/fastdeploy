package entities

import (
	"errors"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
)

// StepDefinition representa la definición de un paso de despliegue (e.g., test, supply).
// Es una Entidad dentro del agregado DeploymentTemplate, identificada por su nombre.
type StepDefinition struct {
	name     string
	commands []vos.CommandDefinition
}

// NewStepDefinition crea una nueva y validada Entidad StepDefinition.
func NewStepDefinition(name string, commands []vos.CommandDefinition) (StepDefinition, error) {
	if name == "" {
		return StepDefinition{}, errors.New("el nombre de la definición de paso no puede estar vacío")
	}
	if len(commands) == 0 {
		return StepDefinition{}, errors.New("un paso debe tener al menos una definición de comando")
	}

	// Creamos una copia de los comandos para asegurar que la entidad sea dueña de sus datos.
	commandsCopy := make([]vos.CommandDefinition, len(commands))
	copy(commandsCopy, commands)

	return StepDefinition{
		name:     name,
		commands: commandsCopy,
	}, nil
}

// Name devuelve el nombre del paso.
func (sd StepDefinition) Name() string {
	return sd.name
}

// Commands devuelve una copia de las definiciones de comando para este paso.
func (sd StepDefinition) Commands() []vos.CommandDefinition {
	commandsCopy := make([]vos.CommandDefinition, len(sd.commands))
	copy(commandsCopy, sd.commands)
	return commandsCopy
}
