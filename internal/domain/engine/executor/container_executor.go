package executor

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

type ContainerExecutor struct {
	baseExecutor    *BaseExecutor
	dockerContainer port.ContainerPort
	variables       *model.StoreEntity
}

func NewContainerExecutor(
	baseExecutor *BaseExecutor,
	dockerContainer port.ContainerPort,
	variables *model.StoreEntity,
) Executor {
	return &ContainerExecutor{
		baseExecutor:    baseExecutor,
		dockerContainer: dockerContainer,
		variables:       variables,
	}
}

func (e *ContainerExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		commitHash := e.variables.Get(constant.KeyCommitHash)
		projectVersion := e.variables.Get(constant.KeyProjectVersion)
		return e.dockerContainer.Start(ctx, commitHash, projectVersion).Error
	})
}
