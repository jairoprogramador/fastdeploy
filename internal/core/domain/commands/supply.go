package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type SupplyCommand struct {
	BaseCommand
	supplyStrategy steps.SupplyStrategy
}

func NewSupplyCommand(strategy steps.SupplyStrategy) Command {
	return &SupplyCommand{
		supplyStrategy: strategy,
	}
}

func (s *SupplyCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: SUPPLY")
	if err := s.supplyStrategy.ExecuteSupply(ctx); err != nil {
		return err
	}
	s.ExecuteNext(ctx)
	return nil
}
