package cmd

import (
	"fmt"
	"os"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type AddControllerSupportFunc func() error
type AddControllerDependencyFunc func(cmd *cobra.Command, args []string) error

func NewAddCmd(addSonarFunc AddControllerSupportFunc, addFortifyFunc AddControllerSupportFunc, addDepFunc AddControllerDependencyFunc) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Agregar recursos al proyecto",
	}

	supportCmd := &cobra.Command{
		Use:   "support",
		Short: "Agrega herramientas de soporte como sonarQube o fortify",
		Run: func(cmd *cobra.Command, args []string) {
			options := []string{"sonarQube", "fortify"}
			prompt := promptui.Select{
				Label: "Seleccione una herramienta de soporte",
				Items: options,
			}
			_, result, errPrompt := prompt.Run()

			if errPrompt != nil {
				return
			}

			switch result {
				case "sonarQube":
					if addSonarFunc == nil {
						fmt.Println("Controlador de comando add support sonarQube no implementado")
						os.Exit(1)
					}
					if err := addSonarFunc(); err != nil {
						os.Exit(1)
					}
				case "fortify":
					if addFortifyFunc == nil {
						fmt.Println("Controlador de comando add support fortify no implementado")
						os.Exit(1)
					}
					if err := addFortifyFunc(); err != nil {
						os.Exit(1)
					}
			}
		},
	}

	dependencyCmd := &cobra.Command{
		Use:   "dependency",
		Short: "Agrega dependencias de proyecto",
		Run: func(cmd *cobra.Command, args []string) {
			if addDepFunc == nil {
				fmt.Println("Controlador de comando add dependency no implementado")
				os.Exit(1)
			}
			if err := addDepFunc(cmd, args); err != nil {
				os.Exit(1)
			}
		},
	}

	addCmd.AddCommand(supportCmd)
	addCmd.AddCommand(dependencyCmd)

	return addCmd
}
