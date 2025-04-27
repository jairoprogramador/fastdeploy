package handler

import (
	"deploy/internal/interface/presenter"
	"deploy/internal/application"
	"github.com/spf13/cobra"
)

func Publish(cmd *cobra.Command, args []string) {
	presenter.ShowStart("Publish")
	message := command.PublishProject()
	if message.Error != nil {
		presenter.ShowError(message.Error)
		return
    }
	presenter.ShowSuccess(message.Message)
}