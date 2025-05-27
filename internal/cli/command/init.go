package cmd

import (
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"github.com/spf13/cobra"
	"os"
)

type InitControllerFunc func() result.DomainResult

func NewInitCmd(initControllerFunc InitControllerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Configura un nuevo proyecto Deploy",
		Run: func(cmd *cobra.Command, args []string) {
			if initControllerFunc != nil {
				if result := initControllerFunc(); !result.IsSuccess() {
					os.Exit(1)
				}
			}
		},
	}
}
