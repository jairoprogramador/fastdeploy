package handler

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"fmt"
)

const (
	errFuncStartNotImplement = "function controller start not implemented"
)

type StartAppFunc func() model.DomainResultEntity

type StartHandler struct {
	startAppFn StartAppFunc
}

func NewStartHandler(startAppFn StartAppFunc) *StartHandler {
	return &StartHandler{
		startAppFn: startAppFn,
	}
}

func (h *StartHandler) Controller() model.DomainResultEntity {
	if h.startAppFn == nil {
		err := fmt.Errorf(errFuncStartNotImplement)
		return model.NewErrorApp(err)
	}

	return h.startAppFn()
}
