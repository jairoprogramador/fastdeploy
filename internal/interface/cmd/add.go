package cmd

import (
	"os"
	"deploy/internal/interface/presenter"
	"github.com/manifoldco/promptui"
	"deploy/internal/interface/handler"
	"github.com/spf13/cobra"
)

func AddCmd() *cobra.Command {
	addCmd :=  &cobra.Command {
		Use:   "add",
		Short: "Agregar recursos",
	}

	supportCmd := &cobra.Command{
		Use:   "support",
		Short: "Agrega herramientas de soporte: sonarQube, fortify",
		Run: func(cmd *cobra.Command, args []string) {
			options := []string{"sonarQube", "fortify"}

			prompt := promptui.Select{
				Label: "Seleccione una herramienta de soporte",
				Items: options,
			}

			_, result, err := prompt.Run()

			if err != nil {
				presenter.ShowError("Add Support", err)
				os.Exit(1)
			}
			
			switch result {
				case "sonarQube":
					handler.AddSupportSonarQube()
				case "fortify":
					handler.AddSupportFortify()
			}
		},
	}

	dependencyCmd := &cobra.Command {
		Use:   "dependency",
		Short: "Agrega dependencias",
		Run: func(cmd *cobra.Command, args []string) {
			handler.AddDependency()
		},
	}

	addCmd.AddCommand(supportCmd)
	addCmd.AddCommand(dependencyCmd)

	return addCmd
}


