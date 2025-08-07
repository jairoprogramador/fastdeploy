package commands

import (
	"fmt"
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

func (d *DeployCommand) Execute() error {
	fmt.Println("Ejecutando el comando: DEPLOY")
	if err := d.deployStrategy.ExecuteDeploy(); err != nil {
		return err
	}
	d.ExecuteNext()
	return nil
}
