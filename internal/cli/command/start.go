package cmd

import (
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"github.com/spf13/cobra"
	"os"
)

type StartControllerFunc func() result.DomainResult

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
