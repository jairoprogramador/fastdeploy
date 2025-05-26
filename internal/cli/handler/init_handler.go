package handler

import (
	"github.com/jairoprogramador/fastdeploy/internal/cli/presenter"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"fmt"
)

const (
	errFuncInitNotImplement = "function controller init not implemented"
)

type InitAppFunc func() model.DomainResultEntity

type InitHandler struct {
	initAppFunc InitAppFunc
}

func NewInitHandler(initAppFunc InitAppFunc) *InitHandler {
	return &InitHandler{
		initAppFunc: initAppFunc,
	}
}

func (h *InitHandler) Controller() model.DomainResultEntity {
	if h.initAppFunc == nil {
		err := fmt.Errorf(errFuncInitNotImplement)
		return model.NewErrorApp(err)
	}

	presenter.ShowBanner()
	return h.initAppFunc()
}
