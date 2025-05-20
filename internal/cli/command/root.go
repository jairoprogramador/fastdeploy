package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var actualRootCmd *cobra.Command

func SetupCommands(rootCmdFn func() *cobra.Command, childCmds ...*cobra.Command) {
	actualRootCmd = rootCmdFn()
	if actualRootCmd == nil {
		fmt.Println("Error: El comando raíz no pudo ser configurado.")
		os.Exit(1)
	}
	actualRootCmd.AddCommand(childCmds...)
}

func Execute() {
	if actualRootCmd == nil {
		fmt.Println("Error: Comandos no configurados. Llamar a SetupCommands primero.")
		os.Exit(1)
	}
	if err := actualRootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

