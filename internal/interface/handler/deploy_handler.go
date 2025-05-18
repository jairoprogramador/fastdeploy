package handler

import (
	"deploy/internal/application"
	"deploy/internal/interface/presenter"
	"os"

	"github.com/spf13/cobra"
)

func Deploy(cmd *cobra.Command, args []string) {
	if cmd.Use != "init" {
		err := application.IsInitialize()
		if err != nil {
			presenter.ShowError("Deploy", err)
			os.Exit(1)
		}
    }
}
