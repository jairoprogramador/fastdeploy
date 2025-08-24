package command

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/chain"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
)

type DeployCommand struct {
	chain.BaseCommandChain
	strategy strategy.StepStrategy
}

func NewDeployCommand(strategy strategy.StepStrategy) chain.CommandChain {
	return &DeployCommand{
		strategy: strategy,
	}
}

func (d *DeployCommand) Execute(ctx deployment.Context) error {
	fmt.Println("Ejecutando el comando: DEPLOY")
	if err := d.strategy.Execute(ctx); err != nil {
		return err
	}
	d.ExecuteNext(ctx)
	return nil
}
