package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type TestCommand struct {
	BaseCommand
	testStrategy strategies.Strategy
}

func NewTestCommand(strategy strategies.Strategy) Command {
	return &TestCommand{
		testStrategy: strategy,
	}
}

func (t *TestCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: TEST")
	if err := t.testStrategy.Execute(ctx); err != nil {
		return err
	}
	t.ExecuteNext(ctx)
	return nil
}
