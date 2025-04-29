package cmd

import (
	"deploy/internal/interface/handler"
	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Configura un nuevo proyecto Deploy",
		Run: func(cmd *cobra.Command, args []string) {
			handler.Initialize()
		},
	}
}

