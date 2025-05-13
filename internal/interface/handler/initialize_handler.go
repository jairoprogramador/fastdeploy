package handler

import (
	"deploy/internal/application"
	"deploy/internal/interface/presenter"
)

func Initialize() {
	presenter.ShowBanner()
	dto := command.InitializeProject()

	if dto.Error != nil {
		presenter.ShowError(dto.Error)
		return
    }
	presenter.ShowSuccess(dto.Message)
}