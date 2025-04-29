package handler

import (
	"deploy/internal/interface/presenter"
	"deploy/internal/application"
)

func Publish() {
	presenter.ShowStart("Publish")
	message := command.PublishProject()
	if message.Error != nil {
		presenter.ShowError(message.Error)
		return
    }
	presenter.ShowSuccess(message.Message)
}