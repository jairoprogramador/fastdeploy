package utils

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
)

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
