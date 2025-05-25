package executor

import (
	"context"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/engine/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/port"
	"deploy/internal/domain/service"
)

// CheckExecutor handles container setup operations
type CheckExecutor struct {
	logger          *logger.Logger
	dockerContainer port.DockerContainer
	router          *service.PathService
	variables       *model.VariableStore
}

// NewSetupExecutor creates a new setup executor instance
func NewCheckExecutor(
	logger *logger.Logger,
	dockerContainer port.DockerContainer,
	variables *model.VariableStore,
	router *service.PathService,
) Executor {
	return &CheckExecutor{
		logger:          logger,
		dockerContainer: dockerContainer,
		router:          router,
		variables:       variables,
	}
}

// Execute checks if a container exists and starts it if needed
func (e *CheckExecutor) Execute(ctx context.Context, step model.Step) error {
	// Check if container exists
	containerExists, err := e.checkContainerExists(ctx)
	if err != nil {
		return err
	}

	if !containerExists {
		e.logger.InfoSystem("no existing container found")
		return nil
	}

	// Container exists, try to start it
	return e.startContainer(ctx)
}

// checkContainerExists verifies if the container already exists
func (e *CheckExecutor) checkContainerExists(ctx context.Context) (bool, error) {
	commitHash := e.variables.Get(constant.VAR_COMMIT_HASH)
	projectVersion := e.variables.Get(constant.VAR_PROJECT_VERSION)

	response := e.dockerContainer.Exists(ctx, commitHash, projectVersion)
	if !response.IsSuccess() {
		e.logger.Error(response.Error)
		return false, response.Error
	}

	return response.Result.(bool), nil
}

// startExistingContainer attempts to start an existing container
func (e *CheckExecutor) startContainer(ctx context.Context) error {
	e.logger.InfoSystem("Existing container detected, attempting to lift container")

	response := e.dockerContainer.Up(ctx)
	if !response.IsSuccess() {
		e.logger.Error(response.Error)
		return response.Error
	}

	// This appears to be a logical error in the original code.
	// It returns an error even when the response is successful.
	// Let's fix it by returning nil instead.
	return nil
}
