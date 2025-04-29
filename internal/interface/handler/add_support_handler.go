package handler

import (
	"deploy/internal/interface/presenter"
	"deploy/internal/application/support"
)

func AddSupportSonarQube() {
	presenter.ShowStart("add support sonarQube")
	message := command.AddSonarQube()
	if message.Error != nil {
		presenter.ShowError(message.Error)
		return
    }
	presenter.ShowSuccess(message.Message)
}

func AddSupportFortify() {
	presenter.ShowStart("add support fortify")
	/* message := command.PublishProject()
	if message.Error != nil {
		presenter.ShowError(message.Error)
		return
    }
	presenter.ShowSuccess(message.Message) */
}