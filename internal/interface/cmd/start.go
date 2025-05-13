package cmd

import (

	"deploy/internal/interface/handler"
	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {
	return &cobra.Command {
		Use:   "start",
		Short: "Publicar aplicación",
		Run:  func(cmd *cobra.Command, args []string) {
			handler.StartPublish()
		},
	}
}


