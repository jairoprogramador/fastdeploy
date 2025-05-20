package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

type InitControllerFunc func() error

func NewInitCmd(runControllerFunc InitControllerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Configura un nuevo proyecto Deploy",
		Run: func(cmd *cobra.Command, args []string) {
			if runControllerFunc == nil {
				fmt.Println("Controlador de comando init no implementado")
				os.Exit(1)
			}
			if err := runControllerFunc(); err != nil {
				os.Exit(1)
			}
		},
	}
}
