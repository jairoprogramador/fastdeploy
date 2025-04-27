package handler

import (
	"deploy/internal/application"
	"deploy/internal/interface/presenter"
	"github.com/spf13/cobra"
)

func Initialize(cmd *cobra.Command, args []string) {
	presenter.ShowBanner()
	message := command.InitializeProject()

	if message.Error != nil {
		presenter.ShowError(message.Error)
		return
    }
	presenter.ShowSuccess(message.Message)
}