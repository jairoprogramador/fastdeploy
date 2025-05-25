package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type StartControllerFunc func() error

func NewStartCmd(startControllerFunc StartControllerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Publicar aplicación",
		Run: func(cmd *cobra.Command, args []string) {
			if startControllerFunc == nil {
				fmt.Println("Controlador de comando start no implementado")
				os.Exit(1)
			}
			if err := startControllerFunc(); err != nil {
				os.Exit(1)
			}
		},
	}
}
