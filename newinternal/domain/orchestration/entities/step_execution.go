package entities

import (
	"errors"
	"fmt"

	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/services"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// StepExecution representa el estado y el resultado de la ejecución de un único paso de despliegue.
// Es una Entidad dentro del agregado Order. Es responsable de gestionar el ciclo de vida
// de sus CommandExecutions y de derivar su propio estado a partir de ellos.
type StepExecution struct {
	name              string
	status            vos.StepStatus
	commandExecutions []*CommandExecution // Usamos punteros para poder modificar su estado.
}

// RehydrateStepExecution reconstruye una entidad StepExecution desde un estado persistido.
func RehydrateStepExecution(name string, status vos.StepStatus, commands []*CommandExecution) *StepExecution {
	return &StepExecution{
		name:              name,
		status:            status,
		commandExecutions: commands,
	}
}

// NewStepExecution crea una nueva instancia de StepExecution.
// Toma una definición de paso (snapshot) y crea la lista de CommandExecutions correspondientes.
func NewStepExecution(def deploymententities.StepDefinition) (*StepExecution, error) {
	if def.Name() == "" {
		return nil, errors.New("la definición del paso debe tener un nombre")
	}

	cmds := def.Commands()
	if len(cmds) == 0 {
		return nil, fmt.Errorf("la definición del paso '%s' no tiene comandos", def.Name())
	}

	commandExecutions := make([]*CommandExecution, 0, len(cmds))
	for _, cmdDef := range cmds {
		cmdExec, err := NewCommandExecution(cmdDef)
		if err != nil {
			return nil, fmt.Errorf("error al crear la ejecución del comando '%s': %w", cmdDef.Name(), err)
		}
		commandExecutions = append(commandExecutions, cmdExec)
	}

	return &StepExecution{
		name:              def.Name(),
		status:            vos.StepStatusPending,
		commandExecutions: commandExecutions,
	}, nil
}

// CompleteCommand encuentra un comando por su nombre, le pasa el resultado de la ejecución,
// y luego actualiza el estado general del paso.
func (se *StepExecution) CompleteCommand(
	commandName, resolvedCmd, log string,
	exitCode int,
	resolver services.VariableResolver,
) error {
	var targetCmd *CommandExecution
	for _, cmd := range se.commandExecutions {
		if cmd.Name() == commandName {
			targetCmd = cmd
			break
		}
	}
	if targetCmd == nil {
		return fmt.Errorf("no se encontró el comando '%s' en el paso '%s'", commandName, se.name)
	}

	err := targetCmd.Execute(resolvedCmd, log, exitCode, resolver)
	if err != nil {
		return err
	}

	// Después de que el comando se completa, actualizamos el estado del paso.
	se.updateStatus()

	return nil
}

// CollectNewVariables recupera las variables generadas por un comando específico.
func (se *StepExecution) CollectNewVariables(commandName string) []vos.Variable {
	for _, cmd := range se.commandExecutions {
		if cmd.Name() == commandName {
			return cmd.OutputVars()
		}
	}
	return nil // No se encontraron variables o no se encontró el comando.
}

// updateStatus recalcula el estado del paso basándose en el estado de sus comandos.
// Esta es una regla de negocio clave del dominio.
func (se *StepExecution) updateStatus() {
	if se.status == vos.StepStatusSkipped {
		return // Si se omite, su estado no cambia.
	}

	hasFailed := false
	allCompleted := true

	for _, cmd := range se.commandExecutions {
		if cmd.Status() == vos.CommandStatusFailed {
			hasFailed = true
			break
		}
		if cmd.Status() == vos.CommandStatusPending {
			allCompleted = false
		}
	}

	if hasFailed {
		se.status = vos.StepStatusFailed
	} else if allCompleted {
		se.status = vos.StepStatusSuccessful
	} else {
		se.status = vos.StepStatusInProgress
	}
}

// Name devuelve el nombre del paso.
func (se *StepExecution) Name() string {
	return se.name
}

// Status devuelve el estado actual del paso.
func (se *StepExecution) Status() vos.StepStatus {
	return se.status
}

// CommandExecutions devuelve una copia de las ejecuciones de comando del paso.
func (se *StepExecution) CommandExecutions() []*CommandExecution {
	cmdsCopy := make([]*CommandExecution, len(se.commandExecutions))
	copy(cmdsCopy, se.commandExecutions)
	return cmdsCopy
}

// Skip marca el paso como omitido.
func (se *StepExecution) Skip() {
	if se.status == vos.StepStatusPending {
		se.status = vos.StepStatusSkipped
	}
}

func (se *StepExecution) MarkAsCached() {
	if se.status == vos.StepStatusPending {
		se.status = vos.StepStatusCached
	}
}
