package commands

import (
	"fmt"
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

func (t *TestCommand) Execute() error {
    fmt.Println("Ejecutando el comando: TEST")
    if err := t.testStrategy.ExecuteTest(); err != nil {
		return err
	}
    t.ExecuteNext()
    return nil
}