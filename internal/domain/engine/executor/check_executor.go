package executor

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
)

type CheckExecutor struct {
	dockerContainer port.DockerContainer
	router          port.PathService
	variables       *model.StoreEntity
}

func NewCheckExecutor(
	dockerContainer port.DockerContainer,
	variables *model.StoreEntity,
	router port.PathService,
) Executor {
	return &CheckExecutor{
		dockerContainer: dockerContainer,
		router:          router,
		variables:       variables,
	}
}

func (e *CheckExecutor) Execute(ctx context.Context, step model.Step) error {
	containerExists, err := e.checkContainerExists(ctx)
	if err != nil {
		return err
	}

	if !containerExists {
		return nil
	}

	return e.startContainer(ctx)
}

func (e *CheckExecutor) checkContainerExists(ctx context.Context) (bool, error) {
	commitHash := e.variables.Get(constant.VAR_COMMIT_HASH)
	projectVersion := e.variables.Get(constant.VAR_PROJECT_VERSION)

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
