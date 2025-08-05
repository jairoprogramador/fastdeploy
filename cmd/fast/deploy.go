package main

import (
	"log"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/spf13/cobra"
)

func NewDeployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Ejecuta el despliegue de la aplicación.",
		Long: `Este comando ejecuta el despliegue de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			command := commands.NewDeployCommand()
			if err := command.Execute(); err != nil {
				log.Fatalf("Error al ejecutar el comando deploy: %v", err)
			}
		},
	}
}
