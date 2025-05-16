package application

import (
	"deploy/internal/domain/engine"
	"deploy/internal/application/dto"
	"deploy/internal/domain/variable"
	"deploy/internal/domain/validator"
	"context"
	"time"
)

func StartDeploy() *dto.ResponseDto {
	storeVariable := variable.GetVariableStore()
	storeService := getStoreService()

	engine := engine.NewEngine(storeVariable, storeService)
	engine.RegisterExecutor(validator.TypeCommand, getCommandExecutor(storeVariable))
	engine.RegisterExecutor(validator.TypeContainer, getContainerExecutor(storeVariable))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()

	deploymentService := getDeploymentService()
	deployment, err := deploymentService.Load()
	if err != nil {
		return dto.GetDtoWithError(err)
	}
	
	if err := engine.Execute(ctx, &deployment); err != nil {
        return dto.GetDtoWithError(err)
    }

	return dto.GetDtoWithMessage("Despliegue completado exitosamente")
}
