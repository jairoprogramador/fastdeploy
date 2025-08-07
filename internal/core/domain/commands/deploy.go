package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/deploy"
)

type DeployCommand struct {
	BaseCommand
	deployStrategy deploy.DeployStrategy
}

func NewDeployCommand(strategy deploy.DeployStrategy) Command {
	return &DeployCommand{
		deployStrategy: strategy,
	}
}

func (d *DeployCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: DEPLOY")
	if err := d.deployStrategy.ExecuteDeploy(); err != nil {
		return err
	}
	/* packageName, err := ctx.Get("package.name")
	if err != nil {
		return err
	}
	packageVersion, err := ctx.Get("package.version")
	if err != nil {
		return err
	}

	fmt.Printf("  Desplegando paquete: %s:%s\n", packageName, packageVersion) */

	d.ExecuteNext(ctx)
	return nil
}
