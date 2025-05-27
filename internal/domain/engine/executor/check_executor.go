package executor

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

type CheckExecutor struct {
	dockerContainer port.DockerContainer
	router          port.PathService
	variables       *entity.StoreEntity
}

func NewCheckExecutor(
	dockerContainer port.DockerContainer,
	variables *entity.StoreEntity,
	router port.PathService,
) Executor {
	return &CheckExecutor{
		dockerContainer: dockerContainer,
		router:          router,
		variables:       variables,
	}
}

func (e *CheckExecutor) Execute(ctx context.Context, step entity.Step) error {
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
