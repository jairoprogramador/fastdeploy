package handler

import (
	"os"
	"fmt"
	"deploy/internal/domain"
	"deploy/internal/application"
	"deploy/internal/interface/presenter"
	"github.com/spf13/cobra"
)

func Deploy(cmd *cobra.Command, args []string) {
	initialized := command.IsInitialized()

	if !initialized && cmd.Use != "init" {
		presenter.ShowError(fmt.Errorf(constants.MessageRunInit))
		os.Exit(1)
    }
}