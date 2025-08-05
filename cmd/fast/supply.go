package main

import (
	"log"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/spf13/cobra"
)

func NewSupplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "supply",
		Short: "Ejecuta el suministro de la aplicación.",
		Long: `Este comando ejecuta el suministro de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			command := commands.NewSupplyCommand()
			if err := command.Execute(); err != nil {
				log.Fatalf("Error al ejecutar el comando supply: %v", err)
			}
		},
	}
}
