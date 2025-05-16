package handler

import (
	"deploy/internal/interface/presenter"
	"deploy/internal/application"
)

func StartPublish() {
	presenter.ShowStart("Start Deploy")
	dto := application.StartDeploy()
	if dto.Error != nil {
		presenter.ShowError(dto.Error)
		return
    }
	presenter.ShowSuccess(dto.Message)
}