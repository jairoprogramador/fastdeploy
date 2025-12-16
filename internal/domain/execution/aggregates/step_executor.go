package aggregates

import (
	"context"
	"fmt"
	"strings"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/entities"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

type StepExecutor struct {
	commandExecutor ports.CommandExecutor
}

func NewStepExecutor(commandExecutor ports.CommandExecutor) *StepExecutor {
	return &StepExecutor{
		commandExecutor: commandExecutor,
	}
}

func (se *StepExecutor) Execute(ctx context.Context, step *entities.Step, initialVars vos.VariableSet) (*vos.ExecutionResult, error) {
	cumulativeLogs := &strings.Builder{}
	cumulativeVars := initialVars.Clone()
	var finalError error
	finalStatus := vos.Success

	for _, command := range step.Commands() {
		cmdResult := se.commandExecutor.Execute(ctx, command, cumulativeVars, step.WorkspaceRoot())

		if cmdResult.Logs != "" {
			cumulativeLogs.WriteString(fmt.Sprintf("--- Log para comando: %s ---\n", command.Name()))
			cumulativeLogs.WriteString(cmdResult.Logs)
			cumulativeLogs.WriteString("\n\n")
		}

		if cmdResult.Error != nil || cmdResult.Status == vos.Failure {
			finalError = fmt.Errorf("el comando '%s' falló", command.Name())
			if cmdResult.Error != nil {
				finalError = fmt.Errorf("el comando '%s' falló: %w", command.Name(), cmdResult.Error)
			}
			finalStatus = vos.Failure
			break
		}

		for key, value := range cmdResult.OutputVars {
			cumulativeVars[key] = value
		}
	}

	return &vos.ExecutionResult{
		Status:     finalStatus,
		Logs:       cumulativeLogs.String(),
		OutputVars: cumulativeVars,
		Error:      finalError,
	}, nil
}
