package handler

import (
	"deploy/internal/cli/presenter"
	"fmt"
	"github.com/spf13/cobra"
)

type AddDependencyAppFunc func(name string, version string) (string, error)

type AddDependencyHandler struct {
	addDependencyFn AddDependencyAppFunc
}

func NewAddDependencyHandler(addDependencyFn AddDependencyAppFunc) *AddDependencyHandler {
	return &AddDependencyHandler{
		addDependencyFn: addDependencyFn,
	}
}

func (h *AddDependencyHandler) Controller(cmd *cobra.Command, args []string) error {
	presenter.ShowStart("AddDependency Command")

	if h.addDependencyFn == nil {
		err := fmt.Errorf("función de aplicación no implementada")
		presenter.ShowError("AddDependencyCmd", err)
		return err
	}

	if len(args) < 2 {
		err := fmt.Errorf("se requieren nombre y versión de la dependencia como argumentos")
		presenter.ShowError("AddDependencyCmd", err)
		return err
	}
	dependencyName := args[0]
	dependencyVersion := args[1]

	message, err := h.addDependencyFn(dependencyName, dependencyVersion)
	if err != nil {
		presenter.ShowError("AddDependencyCmd", err)
		return err
	}
	presenter.ShowSuccess("AddDependencyCmd", message)
	return nil
}
