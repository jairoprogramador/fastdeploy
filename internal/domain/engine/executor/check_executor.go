package executor

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
)

type CheckExecutor struct {
	baseExecutor    *BaseExecutor
	dockerContainer port.ContainerPort
}

func NewCheckExecutor(
	baseExecutor *BaseExecutor,
	dockerContainer port.ContainerPort,
) Executor {
	return &CheckExecutor{
		baseExecutor:    baseExecutor,
		dockerContainer: dockerContainer,
	}
}

func (e *CheckExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		return e.startContainer(ctx)
	})
}

func (e *CheckExecutor) startContainer(ctx context.Context) error {
	response := e.dockerContainer.Up(ctx)
	if !response.IsSuccess() {
		return response.Error
	}
	return nil
}
