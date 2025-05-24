package executor

import (
	"context"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/port"
	"deploy/internal/domain/service/router"
)

type SetupExecutor struct {
	logger          *logger.Logger
	dockerContainer port.DockerContainer
	router          *router.Router
	variables       *model.VariableStore
}

func NewSetupExecutor(
	logger *logger.Logger,
	dockerContainer port.DockerContainer,
	variables *model.VariableStore,
	router *router.Router,
) StepExecutorInterface {
	return &SetupExecutor{
		logger:          logger,
		dockerContainer: dockerContainer,
		router:          router,
		variables:       variables,
	}
}

func (e *SetupExecutor) Execute(ctx context.Context, step model.Step) error {
	response := e.dockerContainer.Exists(ctx, e.variables.Get(constant.VAR_COMMIT_HASH), e.variables.Get(constant.VAR_PROJECT_VERSION))
	if !response.IsSuccess() {
		e.logger.Error(response.Error)
		return response.Error
	}

	exists := response.Result.(bool)
	if exists {
		e.logger.InfoSystem("Existing container detected, attempting to lift container")
		response := e.dockerContainer.Up(ctx)
		if !response.IsSuccess() {
			e.logger.Error(response.Error)
			return response.Error
		}
		message := response.Result.([]string)
		e.logger.ErrorSystemMessage(message[0], response.Error)
		e.logger.Error(response.Error)
		return response.Error
	}
	e.logger.InfoSystem("no existing container found")
	return nil
}
