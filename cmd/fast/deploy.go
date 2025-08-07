package main

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/spf13/cobra"
	"log"
)

func NewDeployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Ejecuta el despliegue de la aplicación.",
		Long:  `Este comando ejecuta el despliegue de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectTechnology := "java" // o "node"

			factory, err := cli.GetStrategyFactory(projectTechnology)
			if err != nil {
				log.Fatalf("Error al obtener la fábrica de estrategias: %v", err)
			}

			testStrategy := factory.CreateTestStrategy()
			supplyStrategy := factory.CreateSupplyStrategy()
			packetStrategy := factory.CreatePackageStrategy()
			deployStrategy := factory.CreateDeployStrategy()

			testCommand := commands.NewTestCommand(testStrategy)
			supplyCommand := commands.NewSupplyCommand(supplyStrategy)
			packageCommand := commands.NewPackageCommand(packetStrategy)
			deployCommand := commands.NewDeployCommand(deployStrategy)

			testCommand.SetNext(supplyCommand)
			supplyCommand.SetNext(packageCommand)
			packageCommand.SetNext(deployCommand)

			pipelineContext := context.NewPipelineContext()

			if err := testCommand.Execute(pipelineContext); err != nil {
				log.Fatalf("Error al ejecutar el comando: %v", err)
			}
		},
	}
}
