package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func SetupCommands(rootCmdFn func() *cobra.Command, childCmds ...*cobra.Command) {
	rootCmd = rootCmdFn()
	if rootCmd == nil {
		fmt.Println("Error: El comando raíz no pudo ser configurado.")
		os.Exit(1)
	}
	rootCmd.AddCommand(childCmds...)
}

func Execute() {
	if rootCmd == nil {
		fmt.Println("Error: Comandos no configurados. Llamar a SetupCommands primero.")
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
