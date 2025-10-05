package entities

import (
	"errors"
	deploymentvos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/services"
	orchestrationvos "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

// CommandExecution representa el estado y el resultado de la ejecución de un único comando.
// Es una Entidad dentro del agregado Order. Su ciclo de vida y estado son gestionados
// por su StepExecution padre y el agregado Order.
type CommandExecution struct {
	name         string
	status       orchestrationvos.CommandStatus
	definition   deploymentvos.CommandDefinition // Snapshot inmutable de la definición.
	resolvedCmd  string                          // Comando con variables resueltas que se ejecutó.
	executionLog string                          // Salida (stdout/stderr) de la ejecución.
	outputVars   []orchestrationvos.Variable                  // Variables extraídas exitosamente de la salida.
}

// NewCommandExecution crea una nueva instancia de CommandExecution.
// Se inicializa en estado "Pending" y toma un snapshot de la definición del comando.
func NewCommandExecution(def deploymentvos.CommandDefinition) (*CommandExecution, error) {
	if def.Name() == "" {
		return nil, errors.New("la definición del comando debe tener un nombre")
	}
	// El struct def se pasa por valor, creando una copia y asegurando el snapshot.
	return &CommandExecution{
		name:       def.Name(),
		status:     orchestrationvos.CommandStatusPending,
		definition: def,
		outputVars: []orchestrationvos.Variable{},
	}, nil
}

// Execute marca la finalización de la ejecución de un comando y procesa el resultado.
// Esta es la función principal que encapsula la lógica de negocio de validación de salida.
func (ce *CommandExecution) Execute(resolvedCmd, log string, exitCode int, resolver services.VariableResolver) (err error) {
	if ce.status != orchestrationvos.CommandStatusPending {
		return errors.New("solo se puede ejecutar un comando que está en estado pendiente")
	}

	ce.resolvedCmd = resolvedCmd
	ce.executionLog = log

	if exitCode != 0 {
		ce.status = orchestrationvos.CommandStatusFailed
		return nil // No es un error del sistema, es un fallo esperado del comando.
	}

	var extractedVars []orchestrationvos.Variable
	for _, probe := range ce.definition.Outputs() {
		variable, match, err := resolver.ExtractVariable(probe, log)
		if err != nil {
			ce.status = orchestrationvos.CommandStatusFailed
			return err
		}

		if !match {
			ce.status = orchestrationvos.CommandStatusFailed
			return err
		}

		if probe.Name() != "" {
			extractedVars = append(extractedVars, variable)
		}
	}

	ce.status = orchestrationvos.CommandStatusSuccessful
	ce.outputVars = extractedVars
	return nil
}

// RehydrateCommandExecution reconstruye una entidad CommandExecution desde un estado persistido.
func RehydrateCommandExecution(name string, status orchestrationvos.CommandStatus, def deploymentvos.CommandDefinition, resolvedCmd, executionLog string, outputVars []orchestrationvos.Variable) *CommandExecution {
	return &CommandExecution{
		name:         name,
		status:       status,
		definition:   def,
		resolvedCmd:  resolvedCmd,
		executionLog: executionLog,
		outputVars:   outputVars,
	}
}

// Name devuelve el nombre del comando.
func (ce *CommandExecution) Name() string {
	return ce.name
}

// Status devuelve el estado actual del comando.
func (ce *CommandExecution) Status() orchestrationvos.CommandStatus {
	return ce.status
}

// ResolvedCmd devuelve el comando que se ejecutó con las variables resueltas.
func (ce *CommandExecution) ResolvedCmd() string {
	return ce.resolvedCmd
}

// ExecutionLog devuelve la salida capturada de la ejecución del comando.
func (ce *CommandExecution) ExecutionLog() string {
	return ce.executionLog
}

// Definition devuelve el snapshot de la definición del comando.
func (ce *CommandExecution) Definition() deploymentvos.CommandDefinition {
	return ce.definition
}

// OutputVars devuelve una copia de las variables que fueron extraídas de la salida.
func (ce *CommandExecution) OutputVars() []orchestrationvos.Variable {
	varsCopy := make([]orchestrationvos.Variable, len(ce.outputVars))
	copy(varsCopy, ce.outputVars)
	return varsCopy
}
