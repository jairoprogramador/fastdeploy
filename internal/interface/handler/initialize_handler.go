package handler

import (
	"deploy/internal/application"
	"deploy/internal/interface/presenter"
)

func Initialize() {
	presenter.ShowBanner()
	logStore := application.Initialize()
	presenter.ShowLogStore(logStore)
}