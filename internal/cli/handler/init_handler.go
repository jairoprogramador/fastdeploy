package handler

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/cli/presenter"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

const (
	errFuncInitNotImplement = "function controller init not implemented"
)

type InitAppFunc func() result.DomainResult

type InitHandler struct {
	initAppFunc InitAppFunc
}

func NewInitHandler(initAppFunc InitAppFunc) *InitHandler {
	return &InitHandler{
		initAppFunc: initAppFunc,
	}
}

func (h *InitHandler) Controller() result.DomainResult {
	if h.initAppFunc == nil {
		err := fmt.Errorf(errFuncInitNotImplement)
		return result.NewErrorApp(err)
	}

	presenter.ShowBanner()
	return h.initAppFunc()
}
