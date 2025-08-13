package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type PackageCommand struct {
	BaseCommand
	strategy strategies.Strategy
}

func NewPackageCommand(strategy strategies.Strategy) Command {
	return &PackageCommand{
		strategy: strategy,
	}
}

func (p *PackageCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: PACKAGE")
	if err := p.strategy.Execute(ctx); err != nil {
		return err
	}
	p.ExecuteNext(ctx)
	return nil
}
