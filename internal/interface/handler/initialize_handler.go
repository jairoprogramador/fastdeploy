package handler

import (
	"deploy/internal/application"
	"deploy/internal/interface/presenter"
)

func Initialize() {
	presenter.ShowBanner()
	message := command.InitializeProject()

	if message.Error != nil {
		presenter.ShowError(message.Error)
		return
    }
	presenter.ShowSuccess(message.Message)
}