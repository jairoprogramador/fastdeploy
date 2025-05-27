package executor

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

type CheckExecutor struct {
	baseExecutor    *BaseExecutor
	dockerContainer port.ContainerPort
	variables       *model.StoreEntity
}

func NewCheckExecutor(
	baseExecutor *BaseExecutor,
	dockerContainer port.ContainerPort,
	variables *model.StoreEntity,
) Executor {
	return &CheckExecutor{
		baseExecutor:    baseExecutor,
		dockerContainer: dockerContainer,
		variables:       variables,
	}
}

func (e *CheckExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		containerExists, err := e.checkContainerExists(ctx)
		if err != nil {
			return err
		}

		if !containerExists {
			return nil
		}

		return e.startContainer(ctx)
	})
}

func (e *CheckExecutor) checkContainerExists(ctx context.Context) (bool, error) {
	commitHash := e.variables.Get(constant.KeyCommitHash)
	projectVersion := e.variables.Get(constant.KeyProjectVersion)

	response := e.dockerContainer.Exists(ctx, commitHash, projectVersion)
	if !response.IsSuccess() {
		return false, response.Error
	}

	return response.Result.(bool), nil
}

func (e *CheckExecutor) startContainer(ctx context.Context) error {
	response := e.dockerContainer.Up(ctx)
	if !response.IsSuccess() {
		return response.Error
	}
	return nil
}
