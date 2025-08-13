package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type PackageCommand struct {
	BaseCommand
	packetStrategy strategies.Strategy
}

func NewPackageCommand(strategy strategies.Strategy) Command {
	return &PackageCommand{
		packetStrategy: strategy,
	}
}

func (p *PackageCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: PACKAGE")
	if err := p.packetStrategy.Execute(ctx); err != nil {
		return err
	}
	p.ExecuteNext(ctx)
	return nil
}
