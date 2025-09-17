package command

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/chain"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
)

type PackageCommand struct {
	chain.BaseCommandChain
	strategy strategy.StepStrategy
}

func NewPackageCommand(strategy strategy.StepStrategy) chain.CommandChain {
	return &PackageCommand{
		strategy: strategy,
	}
}

func (p *PackageCommand) Execute(ctx service.Context) error {
	if err := p.strategy.Execute(ctx); err != nil {
		return err
	}
	return p.ExecuteNext(ctx)
}
