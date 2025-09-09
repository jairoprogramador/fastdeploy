package command

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/chain"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
)

type TestCommand struct {
	chain.BaseCommandChain
	strategy strategy.StepStrategy
}

func NewTestCommand(strategy strategy.StepStrategy) chain.CommandChain {
	return &TestCommand{
		strategy: strategy,
	}
}

func (t *TestCommand) Execute(ctx service.Context) error {
	fmt.Println("Ejecutando el comando: TEST")
	if err := t.strategy.Execute(ctx); err != nil {
		return err
	}
	return t.ExecuteNext(ctx)
}
