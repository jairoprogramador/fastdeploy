package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type TestCommand struct {
	BaseCommand
	strategy strategies.Strategy
}

func NewTestCommand(strategy strategies.Strategy) Command {
	return &TestCommand{
		strategy: strategy,
	}
}

func (t *TestCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: TEST")
	if err := t.strategy.Execute(ctx); err != nil {
		return err
	}
	t.ExecuteNext(ctx)
	return nil
}
