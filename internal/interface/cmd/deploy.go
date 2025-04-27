package cmd

import (
	"deploy/internal/interface/handler"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	if rootCmd != nil {
		return rootCmd
	}

	rootCmd = &cobra.Command{
		Use:   "deploy",
		Short: "CLI para gestionar despliegues de aplicaciones",
		Long:  `Una herramienta de línea de comandos para gestionar el despliegue de aplicaciones en diferentes ambientes con dependencias configurables.`,
		PersistentPreRun: handler.Deploy,
	}

	return rootCmd
}
