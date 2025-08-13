package main

import (
	"fmt"
	"os"

	factory "github.com/jairoprogramador/fastdeploy/internal/adapters/factory/impl"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func NewRootCmd() *cobra.Command {
	if rootCmd != nil {
		return rootCmd
	}

	rootCmd = &cobra.Command{
		Use:   "fd",
		Short: "CLI para gestionar despliegues de aplicaciones",
		Long:  `Una herramienta de línea de comandos para gestionar el despliegue de aplicaciones en diferentes ambientes con dependencias configurables.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() == "init" || cmd.Name() == "config" {
				return
			}

			if !factory.NewInitializeFactory().CreateInitialize().IsInitialized() {
				fmt.Println("El despliegue del proyecto no ha sido inicializado.")
				fmt.Println("Por favor, ejecuta 'init' para comenzar.")
				os.Exit(1)
			}
		},
	}

	return rootCmd
}
