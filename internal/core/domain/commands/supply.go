package commands

import "fmt"

type SupplyCommand struct {}

func NewSupplyCommand() Command {
	return &SupplyCommand{}
}

func (s *SupplyCommand) Execute() error {
    fmt.Println("Ejecutando el comando: SUPPLY")
    return nil
}