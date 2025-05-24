package handler

import (
	"fmt"
	"deploy/internal/cli/presenter"
)

type InitAppFunc func() error

type InitHandler struct {
	initAppFunc InitAppFunc
}

func NewInitHandler(initAppFunc InitAppFunc) *InitHandler {
	return &InitHandler{
		initAppFunc: initAppFunc,
	}
}

func (h *InitHandler) Controller() error {
	if h.initAppFunc == nil {
		err := fmt.Errorf("función de aplicación no implementada")
		presenter.ShowError("InitHandler", err)
		return err
	}

	presenter.ShowBanner()
	return h.initAppFunc()
}
