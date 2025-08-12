package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type DeployCommand struct {
	BaseCommand
	deployStrategy steps.DeployStrategy
}

func NewDeployCommand(strategy steps.DeployStrategy) Command {
	return &DeployCommand{
		deployStrategy: strategy,
	}
}

func (d *DeployCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: DEPLOY")
	if err := d.deployStrategy.ExecuteDeploy(ctx); err != nil {
		return err
	}
	d.ExecuteNext(ctx)
	return nil
}
