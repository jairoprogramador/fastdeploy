package cmd

import (
	"github.com/jairoprogramador/fastdeploy/internal/cli/presenter"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"github.com/spf13/cobra"
	"os"
)

type StartControllerFunc func() result.DomainResult

func NewStartCmd(startControllerFunc StartControllerFunc, fileLogger *logger.FileLogger) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Publicar aplicación",
		Run: func(cmd *cobra.Command, args []string) {
			if startControllerFunc != nil {
				result := startControllerFunc()
				presenter.Show(result, fileLogger)
				if !result.IsSuccess() {
					os.Exit(1)
				}
			}
		},
	}
}
