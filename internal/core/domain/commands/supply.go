package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type SupplyCommand struct {
	BaseCommand
	strategy strategies.Strategy
}

func NewSupplyCommand(strategy strategies.Strategy) Command {
	return &SupplyCommand{
		strategy: strategy,
	}
}

func (s *SupplyCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: SUPPLY")
	if err := s.strategy.Execute(ctx); err != nil {
		return err
	}
	s.ExecuteNext(ctx)
	return nil
}
