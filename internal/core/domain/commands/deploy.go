package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type DeployCommand struct {
	BaseCommand
	strategy strategies.Strategy
}

func NewDeployCommand(strategy strategies.Strategy) Command {
	return &DeployCommand{
		strategy: strategy,
	}
}

func (d *DeployCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: DEPLOY")
	if err := d.strategy.Execute(ctx); err != nil {
		return err
	}
	d.ExecuteNext(ctx)
	return nil
}
