package command

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/chain"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
)

type SupplyCommand struct {
	chain.BaseCommandChain
	strategy strategy.StepStrategy
}

func NewSupplyCommand(strategy strategy.StepStrategy) chain.CommandChain {
	return &SupplyCommand{
		strategy: strategy,
	}
}

func (s *SupplyCommand) Execute(ctx service.Context) error {
	fmt.Println("Ejecutando el comando: SUPPLY")
	if err := s.strategy.Execute(ctx); err != nil {
		return err
	}
	return s.ExecuteNext(ctx)
}
