package command

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
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

func (d *DeployCommand) Execute(ctx service.Context) error {
	if err := d.strategy.Execute(ctx); err != nil {
		return err
	}
	return d.ExecuteNext(ctx)
}
