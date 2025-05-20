package executor

import (
	"context"
	"fmt"
	"deploy/internal/domain/model"
	"deploy/internal/domain/router"
	"deploy/internal/domain/service"
)

type SetupExecutor struct {
	dockerService service.DockerServiceInterface
	router        *router.Router
	logStore      *model.LogStore
	variables     *model.VariableStore
}

func NewSetupExecutor(
	dockerService service.DockerServiceInterface,
	variables *model.VariableStore,
	router *router.Router,
	logStore *model.LogStore,
) StepExecutorInterface {
	return &SetupExecutor{
		dockerService: dockerService,
		router:        router,
		logStore:      logStore,
		variables:     variables,
	}
}

func (e *SetupExecutor) Execute(ctx context.Context, step model.Step) (string, error) {
	exists, err := e.dockerService.ExistsContainer(ctx, e.variables)
	if err != nil {
		e.logStore.AddError(fmt.Errorf("error verificando existencia de contenedor en SetupContainersExecutor: %v", err))
		return "", err
	}

	if exists {
		e.logStore.AddMessage("Contenedor existente detectado, intentando levantar con Docker Compose.")
		pathDockerCompose := e.router.GetFullPathDockerCompose()
		message, err := e.dockerService.DockerComposeUp(ctx, pathDockerCompose, e.variables)
		if err == nil {
			e.logStore.AddMessage(message)
			return message, nil
		}
		e.logStore.AddError(fmt.Errorf("error en SetupContainersExecutor intentando levantar compose: %v", err))
		return "", err
	}
	e.logStore.AddMessage("No se detectó contenedor existente para pre-configurar.")
	return "No se detectó contenedor existente, no se requiere acción de pre-configuración.", nil
}
