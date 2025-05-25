package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var deployCmdInstance *cobra.Command

type DeployControllerFunc func() error

func GetDeployCmd(deployControllerFunc DeployControllerFunc) *cobra.Command {
	if deployCmdInstance != nil {
		return deployCmdInstance
	}

	deployCmdInstance = &cobra.Command{
		Use:   "deploy",
		Short: "CLI para gestionar despliegues de aplicaciones",
		Long:  `Una herramienta de línea de comandos para gestionar el despliegue de aplicaciones en diferentes ambientes con dependencias configurables.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if deployControllerFunc == nil {
				fmt.Println("Controlador de comando deploy no implementado")
				os.Exit(1)
			}
			if err := deployControllerFunc(); err != nil {
				os.Exit(1)
			}
		},
	}

	return deployCmdInstance
}
