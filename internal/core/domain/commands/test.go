package commands

import "fmt"

type TestCommand struct {}

func NewTestCommand() Command {
	return &TestCommand{}
}

func (t *TestCommand) Execute() error {
    fmt.Println("Ejecutando el comando: TEST")
    return nil
}