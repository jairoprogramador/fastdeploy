package executor

import (
	"context"
	"deploy/internal/domain/engine/model"
	"deploy/internal/domain/port"
)

// ContainerExecutor handles Docker container operations
type ContainerExecutor struct {
	baseExecutor    *StepExecutor
	variables       *model.VariableStore
	dockerContainer port.DockerContainer
}

// NewContainerExecutor creates a new container executor instance
func NewContainerExecutor(
	baseExecutor *StepExecutor,
	variables *model.VariableStore,
	dockerContainer port.DockerContainer,
) Executor {
	return &ContainerExecutor{
		baseExecutor:    baseExecutor,
		variables:       variables,
		dockerContainer: dockerContainer,
	}
}

// Execute runs the container operation defined in the step
func (e *ContainerExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		// Set up variable scope for this execution
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()

		// Start the container
		return e.dockerContainer.Start(ctx)
	})
}
