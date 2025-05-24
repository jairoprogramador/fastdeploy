package handler

import (
	"deploy/internal/cli/presenter"
	"deploy/internal/domain/model"
	"fmt"
)

type StartAppFunc func(project *model.ProjectEntity) error

type StartHandler struct {
	startAppFn  StartAppFunc
	isInitAppFn IsInitAppFunc
}

func NewStartHandler(startAppFn StartAppFunc, isInitAppFn IsInitAppFunc) *StartHandler {
	return &StartHandler{
		startAppFn:  startAppFn,
		isInitAppFn: isInitAppFn,
	}
}

func (h *StartHandler) Controller() error {
	if h.startAppFn == nil {
		err := fmt.Errorf("start: función de aplicación no implementada")
		presenter.ShowError("StartHandler", err)
		return err
	}

	if h.isInitAppFn == nil {
		err := fmt.Errorf("isInit: función de aplicación no implementada")
		presenter.ShowError("StartHandler", err)
		return err
	}

	presenter.ShowStart("Start Deploy Command")

	project, err := h.isInitAppFn()
	if err != nil {
		presenter.ShowError("StartCmd (IsInitialize)", err)
		return err
	}
	if project == nil {
		err = fmt.Errorf("el proyecto no pudo ser cargado o no está inicializado")
		presenter.ShowError("StartCmd (IsInitialize)", err)
		return err
	}

	return h.startAppFn(project)
}
