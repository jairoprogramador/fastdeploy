package main

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/spf13/cobra"
	"log"
)

func NewPackageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "package",
		Short: "Ejecuta el empaquetado de la aplicación.",
		Long:  `Este comando ejecuta el empaquetado de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectTechnology := "java" // o "node"

			factory, err := cli.GetStrategyFactory(projectTechnology)
			if err != nil {
				log.Fatalf("Error al obtener la fábrica de estrategias: %v", err)
			}

			testStrategy := factory.CreateTestStrategy()
			supplyStrategy := factory.CreateSupplyStrategy()
			packetStrategy := factory.CreatePackageStrategy()

			testCommand := commands.NewTestCommand(testStrategy)
			supplyCommand := commands.NewSupplyCommand(supplyStrategy)
			packageCommand := commands.NewPackageCommand(packetStrategy)

			testCommand.SetNext(supplyCommand)
			supplyCommand.SetNext(packageCommand)

			if err := testCommand.Execute(); err != nil {
				log.Fatalf("Error al ejecutar el comando package: %v", err)
			}
		},
	}
}
