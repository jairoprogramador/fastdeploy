package handler

import (
	"deploy/internal/application"
	"deploy/internal/interface/presenter"
)

func Initialize() {
	presenter.ShowBanner()
	dto := application.Initialize()

	if dto.Error != nil {
		presenter.ShowError(dto.Error)
		return
    }
	presenter.ShowSuccess(dto.Message)
}