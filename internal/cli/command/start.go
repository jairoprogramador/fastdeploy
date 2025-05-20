package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

type StartControllerFunc func() error

func NewStartCmd(runControllerFunc StartControllerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Publicar aplicación",
		Run: func(cmd *cobra.Command, args []string) {
			if runControllerFunc == nil {
				fmt.Println("Controlador de comando start no implementado")
				os.Exit(1)
			}
			if err := runControllerFunc(); err != nil {
				os.Exit(1)
			}
		},
	}
}
