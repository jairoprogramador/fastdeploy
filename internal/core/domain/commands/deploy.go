package commands

import "fmt"

type DeployCommand struct {}

func NewDeployCommand() Command {
	return &DeployCommand{}
}

func (d *DeployCommand) Execute() error {
    fmt.Println("Ejecutando el comando: DEPLOY")
    return nil
}