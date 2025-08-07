package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/supply"
)

type SupplyCommand struct {
	BaseCommand
	supplyStrategy supply.SupplyStrategy
}

func NewSupplyCommand(strategy supply.SupplyStrategy) Command {
	return &SupplyCommand{
		supplyStrategy: strategy,
	}
}

func (s *SupplyCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: SUPPLY")
	if err := s.supplyStrategy.ExecuteSupply(); err != nil {
		return err
	}
	s.ExecuteNext(ctx)
	return nil
}
