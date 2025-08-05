package main

import (
	"log"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/spf13/cobra"
)

func NewTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Ejecuta las pruebas de calidad del software.",
		Long: `Este comando ejecuta pruebas unitarias, de integración, escaneos de seguridad y otros análisis estáticos
			para asegurar la calidad del código.`,
		Run: func(cmd *cobra.Command, args []string) {
			command := commands.NewTestCommand()
			if err := command.Execute(); err != nil {
				log.Fatalf("Error al ejecutar el comando test: %v", err)
			}
		},
	}
}
