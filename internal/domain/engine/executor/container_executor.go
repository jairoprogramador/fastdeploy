package executor

import (
	"context"
	"deploy/internal/domain/model"
	"deploy/internal/domain/port"
)

type ContainerExecutor struct {
	baseExecutor    *BaseExecutor
	variables       *model.VariableStore
	dockerContainer port.DockerContainer
}

func NewContainerExecutor(
	baseExecutor *BaseExecutor,
	variables *model.VariableStore,
	dockerContainer port.DockerContainer,
) StepExecutorInterface {
	return &ContainerExecutor{
		baseExecutor:    baseExecutor,
		variables:       variables,
		dockerContainer: dockerContainer,
	}
}

func (e *ContainerExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()
		return e.dockerContainer.Start(ctx)
	})
}
