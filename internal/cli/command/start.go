package cmd

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/spf13/cobra"
	"os"
)

type StartControllerFunc func() model.DomainResultEntity

func NewStartCmd(startControllerFunc StartControllerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Publicar aplicación",
		Run: func(cmd *cobra.Command, args []string) {
			if startControllerFunc != nil {
				if result := startControllerFunc(); !result.IsSuccess() {
					os.Exit(1)
				}
			}
		},
	}
}
