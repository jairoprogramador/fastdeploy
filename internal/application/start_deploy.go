package application

import (
	"deploy/internal/domain/engine"
	"deploy/internal/domain/model"
	"deploy/internal/domain/validator"
	"context"
	"time"
)

func StartDeploy() *model.LogStore {
	logStore := model.NewLogStore("start deploy")
	storeVariable := model.GetVariableStore()
	storeService := getStoreService()

	engine := engine.NewEngine(storeVariable, storeService)
	engine.RegisterExecutor(validator.TypeCommand, getCommandExecutor(storeVariable))
	engine.RegisterExecutor(validator.TypeContainer, getContainerExecutor(storeVariable))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()

	deploymentService := getDeploymentService()
	deployment, err := deploymentService.Load()
	if err != nil {
		logStore.AddError(err)
	}
	
	if err := engine.Execute(ctx, &deployment); err != nil {
        logStore.AddError(err)
    }

	logStore.FinishSteps()
	return logStore
}
