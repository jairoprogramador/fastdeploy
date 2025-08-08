package main

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/spf13/cobra"
	"log"
)

func NewPackageCmd() *cobra.Command {
	skippableSteps := []string{constants.StepTest, constants.StepSupply}

	cmd := &cobra.Command{
		Use:   "package",
		Short: "Ejecuta el empaquetado de la aplicación.",
		Long:  `Este comando ejecuta el empaquetado de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectTechnology := "java" // o "node"

			factory, err := cli.GetStrategyFactory(projectTechnology)
			if err != nil {
				log.Fatalf("Error al obtener la fábrica de estrategias: %v", err)
			}

			allCommands := map[string]commands.Command{
				constants.StepTest:    commands.NewTestCommand(factory.CreateTestStrategy()),
				constants.StepSupply:  commands.NewSupplyCommand(factory.CreateSupplyStrategy()),
				constants.StepPackage: commands.NewPackageCommand(factory.CreatePackageStrategy()),
			}

			skipFlags := cli.GetSkipFlags(cmd, skippableSteps)

			executionOrder := []string{constants.StepTest, constants.StepSupply, constants.StepPackage}

			firstCommand, err := cli.BuildDynamicChain(allCommands, skipFlags, executionOrder)
			if err != nil {
				log.Fatalf("Error al construir la cadena de comandos: %v", err)
			}

			if firstCommand != nil {
				pipelineContext := context.NewPipelineContext()

				if err := firstCommand.Execute(pipelineContext); err != nil {
					log.Fatalf("Error al ejecutar el comando: %v", err)
				}
			} else {
				fmt.Println("No se seleccionaron pasos para ejecutar. Saliendo...")
			}
		},
	}
	cli.AddSkipFlags(cmd, skippableSteps)
	return cmd
}
