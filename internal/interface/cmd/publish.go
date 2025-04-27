package cmd

import (

	"deploy/internal/interface/handler"
	"github.com/spf13/cobra"
)

func PublishCmd() *cobra.Command {
	return &cobra.Command {
		Use:   "publish",
		Short: "Publicar aplicación",
		Run:   handler.Publish,
	}
}


