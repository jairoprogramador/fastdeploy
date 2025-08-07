package main

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/spf13/cobra"
	"log"
)

func NewSupplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "supply",
		Short: "Ejecuta el suministro de la aplicación.",
		Long:  `Este comando ejecuta el suministro de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectTechnology := "java" // o "node"

			factory, err := cli.GetStrategyFactory(projectTechnology)
			if err != nil {
				log.Fatalf("Error al obtener la fábrica de estrategias: %v", err)
			}

			testStrategy := factory.CreateTestStrategy()
			supplyStrategy := factory.CreateSupplyStrategy()

			testCommand := commands.NewTestCommand(testStrategy)
			supplyCommand := commands.NewSupplyCommand(supplyStrategy)

			testCommand.SetNext(supplyCommand)

			pipelineContext := context.NewPipelineContext()

			if err := testCommand.Execute(pipelineContext); err != nil {
				log.Fatalf("Error al ejecutar el comando supply: %v", err)
			}
		},
	}
}
