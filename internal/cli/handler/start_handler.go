package handler

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

const (
	errFuncStartNotImplement = "function controller start not implemented"
)

type StartAppFunc func() result.DomainResult

type StartHandler struct {
	startAppFn StartAppFunc
}

func NewStartHandler(startAppFn StartAppFunc) *StartHandler {
	return &StartHandler{
		startAppFn: startAppFn,
	}
}

func (h *StartHandler) Controller() result.DomainResult {
	if h.startAppFn == nil {
		return result.NewErrorApp(fmt.Errorf(errFuncStartNotImplement))
	}
	return h.startAppFn()
}
