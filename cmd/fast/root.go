package main

import (
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
	}

	return rootCmd
}
