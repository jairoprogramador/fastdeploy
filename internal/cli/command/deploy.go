package cmd

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/spf13/cobra"
	"os"
)

var deployCmdInstance *cobra.Command

type DeployControllerFunc func() model.DomainResultEntity

func GetDeployCmd(deployControllerFunc DeployControllerFunc) *cobra.Command {
	if deployCmdInstance != nil {
		return deployCmdInstance
	}

	deployCmdInstance = &cobra.Command{
		Use:   "deploy",
		Short: "CLI para gestionar despliegues de aplicaciones",
		Long:  `Una herramienta de línea de comandos para gestionar el despliegue de aplicaciones en diferentes ambientes con dependencias configurables.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if deployControllerFunc != nil {
				if result := deployControllerFunc(); !result.IsSuccess() {
					os.Exit(1)
				}
			}
		},
	}

	return deployCmdInstance
}
