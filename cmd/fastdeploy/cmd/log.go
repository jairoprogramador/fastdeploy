package cmd

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/factory"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show the detailed log of the last execution",
	Long:  `Reads and displays the most recent log file from the .fastdeploy/logs directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		factory := factory.NewFactory()
		logService := factory.BuildLogService()

		err := logService.ShowLog()
		if err != nil {
			return err
		}

		return nil
	},
}
