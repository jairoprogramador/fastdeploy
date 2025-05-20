package handler

import (
	"fmt"
	"deploy/internal/domain/model"
	"deploy/internal/cli/presenter"
)

type InitAppFunc func(logStore *model.LogStore) *model.LogStore

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
	initLogStore := model.NewLogStore("InitializeOperationCmd")
	finalLogStore := h.initAppFunc(initLogStore)
	presenter.ShowLogStore(finalLogStore)
	return finalLogStore.GetError()
}
