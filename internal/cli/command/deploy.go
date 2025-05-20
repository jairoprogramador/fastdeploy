package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

var deployCmdInstance *cobra.Command

type DeployControllerFunc func() error

func GetDeployCmd(runControllerFunc DeployControllerFunc) *cobra.Command {
	if deployCmdInstance != nil {
		return deployCmdInstance
	}

	deployCmdInstance = &cobra.Command{
		Use:   "deploy",
		Short: "CLI para gestionar despliegues de aplicaciones",
		Long:  `Una herramienta de línea de comandos para gestionar el despliegue de aplicaciones en diferentes ambientes con dependencias configurables.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if runControllerFunc == nil {
				fmt.Println("Controlador de comando deploy no implementado")
				os.Exit(1)
			}
			if err := runControllerFunc(); err != nil {
				os.Exit(1)
			}
		},
	}

	return deployCmdInstance
}
