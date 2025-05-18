package handler

import (
	"deploy/internal/interface/presenter"
	"deploy/internal/application"
)

func StartPublish() {
	presenter.ShowStart("Start Deploy")
	logStore := application.StartDeploy()
	presenter.ShowLogStore(logStore)
}