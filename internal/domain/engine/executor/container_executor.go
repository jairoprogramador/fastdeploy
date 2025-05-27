package executor

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
)

type ContainerExecutor struct {
	baseExecutor    *BaseExecutor
	dockerContainer port.ContainerPort
}

func NewContainerExecutor(
	baseExecutor *BaseExecutor,
	dockerContainer port.ContainerPort,
) Executor {
	return &ContainerExecutor{
		baseExecutor:    baseExecutor,
		dockerContainer: dockerContainer,
	}
}

func (e *ContainerExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		return e.dockerContainer.Start(ctx).Error
	})
}
