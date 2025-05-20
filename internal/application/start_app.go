package application

import (
	"context"
	"time"
	"deploy/internal/domain/engine"
	"deploy/internal/domain/model"
	"deploy/internal/domain/service"
)

func StartDeploy(
	engineInstance *engine.Engine,
	deploymentService service.DeploymentServiceInterface,
	project *model.Project,
) *model.LogStore {
	logStore := model.NewLogStore("start deploy")

	deployment, err := deploymentService.Load()
	if err != nil {
		logStore.AddError(err)
	}else{
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := engineInstance.Execute(ctx, deployment, project); err != nil {
			logStore.AddError(err)
		}
	}

	logStore.FinishSteps()
	return logStore
}
