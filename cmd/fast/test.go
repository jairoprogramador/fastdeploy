package main

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/spf13/cobra"
	"log"
)

func NewTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Ejecuta las pruebas de calidad del software.",
		Long: `Este comando ejecuta pruebas unitarias, de integración, escaneos de seguridad y otros análisis estáticos
			para asegurar la calidad del código.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectTechnology := "java" // o "node"

			factory, err := cli.GetStrategyFactory(projectTechnology)
			if err != nil {
				log.Fatalf("Error al obtener la fábrica de estrategias: %v", err)
			}

			testStrategy := factory.CreateTestStrategy()

			testCommand := commands.NewTestCommand(testStrategy)

			pipelineContext := context.NewPipelineContext()

			if err := testCommand.Execute(pipelineContext); err != nil {
				log.Fatalf("Error al ejecutar el comando test: %v", err)
			}
		},
	}
}
