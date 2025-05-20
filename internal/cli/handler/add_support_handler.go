package handler

import (
	"deploy/internal/cli/presenter"
	"fmt"
)

type AddSonarQubeAppFunc func() (string, error)

type AddFortifyAppFunc func() (string, error)

type AddSupportHandler struct {
	addSonarQubeFn AddSonarQubeAppFunc
	addFortifyFn   AddFortifyAppFunc
}

func NewAddSupportHandler(addSonarQubeFn AddSonarQubeAppFunc, addFortifyFn AddFortifyAppFunc) *AddSupportHandler {
	return &AddSupportHandler{
		addSonarQubeFn: addSonarQubeFn,
		addFortifyFn:   addFortifyFn,
	}
}

func (h *AddSupportHandler) ControllerSonarQube() error {
	if h.addSonarQubeFn == nil {
		err := fmt.Errorf("función de aplicación no implementada")
		presenter.ShowError("AddSupportSonarQubeCmd", err)
		return err
	}

	presenter.ShowStart("add support sonarQube")
	message, err := h.addSonarQubeFn()
	if err != nil {
		presenter.ShowError("AddSupportSonarQubeCmd", err)
		return err
	}
	presenter.ShowSuccess("AddSupportSonarQubeCmd", message)
	return nil
}

func (h *AddSupportHandler) ControllerFortify() error {
	if h.addFortifyFn == nil {
		err := fmt.Errorf("función de aplicación no implementada")
		presenter.ShowError("AddSupportFortifyCmd", err)
		return err
	}

	presenter.ShowStart("add support fortify")
	message, err := h.addFortifyFn()
	if err != nil {
		presenter.ShowError("AddSupportFortifyCmd", err)
		return err
	}
	presenter.ShowSuccess("AddSupportFortifyCmd", message)
	return nil
}
