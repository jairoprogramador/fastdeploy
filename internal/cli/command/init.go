package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type InitControllerFunc func() error

func NewInitCmd(initControllerFunc InitControllerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Configura un nuevo proyecto Deploy",
		Run: func(cmd *cobra.Command, args []string) {
			if initControllerFunc == nil {
				fmt.Println("Controlador de comando init no implementado")
				os.Exit(1)
			}
			if err := initControllerFunc(); err != nil {
				os.Exit(1)
			}
		},
	}
}
