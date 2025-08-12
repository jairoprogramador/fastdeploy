package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type PackageCommand struct {
	BaseCommand
	packetStrategy steps.PacketStrategy
}

func NewPackageCommand(strategy steps.PacketStrategy) Command {
	return &PackageCommand{
		packetStrategy: strategy,
	}
}

func (p *PackageCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: PACKAGE")
	if err := p.packetStrategy.ExecutePacket(ctx); err != nil {
		return err
	}
	p.ExecuteNext(ctx)
	return nil
}
