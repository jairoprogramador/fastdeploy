package utils

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/spf13/cobra"
)

func AddSkipFlags(cmd *cobra.Command, steps []string) {
	for _, step := range steps {
		flagName := fmt.Sprintf("%s%s", constants.SkipFlagPrefix, step)
		shortFlag := string(step[0])
		description := fmt.Sprintf("Omite el paso de %s", step)
		cmd.Flags().BoolP(flagName, shortFlag, false, description)
	}
}

func GetSkipFlags(cmd *cobra.Command, steps []string) map[string]bool {
	skipFlags := make(map[string]bool)
	for _, step := range steps {
		flagName := fmt.Sprintf("%s%s", constants.SkipFlagPrefix, step)
		value, err := cmd.Flags().GetBool(flagName)
		if err != nil {
			skipFlags[step] = false
		} else {
			skipFlags[step] = value
		}
	}
	return skipFlags
}
