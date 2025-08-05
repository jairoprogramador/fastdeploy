package commands

import "fmt"

type PackageCommand struct {}

func NewPackageCommand() Command {
	return &PackageCommand{}
}

func (p *PackageCommand) Execute() error {
    fmt.Println("Ejecutando el comando: PACKAGE")
    return nil
}