package cli

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
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

func BuildDynamicChain(
	allCommands map[string]commands.Command,
	skipFlags map[string]bool,
	executionOrder []string,
) (commands.Command, error) {
	var firstCommand commands.Command
	var lastCommand commands.Command

	addCommand := func(c commands.Command) {
		if firstCommand == nil {
			firstCommand = c
			lastCommand = c
		} else {
			lastCommand.SetNext(c)
			lastCommand = c
		}
	}

	for _, stepName := range executionOrder {
		if !skipFlags[stepName] {
			cmd, ok := allCommands[stepName]
			if !ok {
				return nil, fmt.Errorf("comando no encontrado en el mapa: %s", stepName)
			}
			addCommand(cmd)
		}
	}

	return firstCommand, nil
}
