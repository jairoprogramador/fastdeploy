package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type TestCommand struct {
	BaseCommand
	testStrategy steps.TestStrategy
}

func NewTestCommand(strategy steps.TestStrategy) Command {
	return &TestCommand{
		testStrategy: strategy,
	}
}

func (t *TestCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: TEST")
	if err := t.testStrategy.ExecuteTest(ctx); err != nil {
		return err
	}
	t.ExecuteNext(ctx)
	return nil
}
