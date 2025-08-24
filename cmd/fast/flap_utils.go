package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	SkipFlagPrefix = "skip-"
)

func AddSkipFlags(cmd *cobra.Command, steps []string) {
	for _, step := range steps {
		flagName := fmt.Sprintf("%s%s", SkipFlagPrefix, step)
		shortFlag := string(step[0])
		description := fmt.Sprintf("Omite el paso de %s", step)
		cmd.Flags().BoolP(flagName, shortFlag, false, description)
	}
}

func GetSkipSteps(cmd *cobra.Command, steps []string) []string {
	var skipSteps []string
	for _, step := range steps {
		flagName := fmt.Sprintf("%s%s", SkipFlagPrefix, step)
		isFlag, _ := cmd.Flags().GetBool(flagName)
		if isFlag {
			skipSteps = append(skipSteps, step)
		}
	}
	return skipSteps
}
