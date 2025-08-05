package main

import (
	"log"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/spf13/cobra"
)

func NewPackageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "package",
		Short: "Ejecuta el empaquetado de la aplicación.",
		Long: `Este comando ejecuta el empaquetado de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			command := commands.NewPackageCommand()
			if err := command.Execute(); err != nil {
				log.Fatalf("Error al ejecutar el comando package: %v", err)
			}
		},
	}
}
