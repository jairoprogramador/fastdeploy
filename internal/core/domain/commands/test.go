package commands

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/test"
)

type TestCommand struct {
	BaseCommand
	testStrategy test.TestStrategy
}

func NewTestCommand(strategy test.TestStrategy) Command {
	return &TestCommand{
		testStrategy: strategy,
	}
}

func (t *TestCommand) Execute(ctx context.Context) error {
	fmt.Println("Ejecutando el comando: TEST")
	if err := t.testStrategy.ExecuteTest(); err != nil {
		return err
	}
	//ctx.Set("package.name", "packageName test")
	//ctx.Set("package.version", "packageVersion test")
	t.ExecuteNext(ctx)
	return nil
}
