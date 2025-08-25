package command

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
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

func (p *PackageCommand) Execute(ctx deployment.Context) error {
	fmt.Println("Ejecutando el comando: PACKAGE")
	if err := p.strategy.Execute(ctx); err != nil {
		return err
	}
	return p.ExecuteNext(ctx)
}
