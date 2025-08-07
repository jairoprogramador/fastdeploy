package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/packet"
)

type PackageCommand struct {
	BaseCommand
	packetStrategy packet.PacketStrategy
}

func NewPackageCommand(strategy packet.PacketStrategy) Command {
	return &PackageCommand{
		packetStrategy: strategy,
	}
}

func (p *PackageCommand) Execute() error {
	fmt.Println("Ejecutando el comando: PACKAGE")
	if err := p.packetStrategy.ExecutePacket(); err != nil {
		return err
	}
	p.ExecuteNext()
	return nil
}
