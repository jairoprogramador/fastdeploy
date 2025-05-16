package handler

import (
	"os"
	"deploy/internal/application"
	"deploy/internal/interface/presenter"
	"github.com/spf13/cobra"
)

func Deploy(cmd *cobra.Command, args []string) {
	if cmd.Use != "init" {
		dto := application.IsInitialize()
		if dto.Error != nil {
			presenter.ShowError(dto.Error)
			os.Exit(1)
		}
    }
}
