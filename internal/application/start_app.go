package application

import (
	"context"
	"deploy/internal/domain/engine"
	"deploy/internal/domain/model"
	"deploy/internal/domain/service"
	"time"
)

func StartDeploy(
	engineInstance *engine.Engine,
	deploymentService service.DeploymentLoader,
	project *model.ProjectEntity,
) error {
	deployment, err := deploymentService.Load()
	if err != nil {
		return err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := engineInstance.Execute(ctx, deployment, project); err != nil {
			return err
		}
	}
	return nil
}
