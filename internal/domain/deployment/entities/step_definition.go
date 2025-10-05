package entities

import (
	"errors"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
)

// StepDefinition representa la definición de un paso de despliegue (e.g., test, supply).
// Es una Entidad dentro del agregado DeploymentTemplate, identificada por su nombre.
type StepDefinition struct {
	name              string
	verificationTypes []vos.VerificationType // <-- AÑADIDO
	commands          []vos.CommandDefinition
}

// NewStepDefinition crea una nueva y validada Entidad StepDefinition.
func NewStepDefinition(name string, verifications []vos.VerificationType, commands []vos.CommandDefinition) (StepDefinition, error) {
	if name == "" {
		return StepDefinition{}, errors.New("el nombre de la definición de paso no puede estar vacío")
	}

	verificationsCopy := make([]vos.VerificationType, len(verifications))
	copy(verificationsCopy, verifications)

	commandsCopy := make([]vos.CommandDefinition, len(commands))
	copy(commandsCopy, commands)

	return StepDefinition{
		name:              name,
		verificationTypes: verificationsCopy, // <-- AÑADIDO
		commands:          commandsCopy,
	}, nil
}

// Name devuelve el nombre del paso.
func (sd StepDefinition) Name() string {
	return sd.name
}

// VerificationTypes devuelve los tipos de verificación para el paso.
func (sd StepDefinition) VerificationTypes() []vos.VerificationType {
	verificationsCopy := make([]vos.VerificationType, len(sd.verificationTypes))
	copy(verificationsCopy, sd.verificationTypes)
	return verificationsCopy
}

// Commands devuelve una copia de las definiciones de comando para este paso.
func (sd StepDefinition) Commands() []vos.CommandDefinition {
	commandsCopy := make([]vos.CommandDefinition, len(sd.commands))
	copy(commandsCopy, sd.commands)
	return commandsCopy
}
