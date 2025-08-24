package service

import (
	"fmt"
	"os"
	"os/exec"
)

type ExecutorCmd interface {
	Execute(yamlFilePath string) error
}

type CommandExecutor struct{}

func NewCommandExecutor() ExecutorCmd {
	return &CommandExecutor{}
}

func (e *CommandExecutor) Execute(yamlFilePath string) error {
	listCmd, err := Load(yamlFilePath)
	if err != nil {
		return err
	}

	for _, cmdDef := range listCmd.Commands {
		fmt.Printf("    -> %s\n", cmdDef.Name)

		projectDir := "."

		cmd := exec.Command("sh", "-c", cmdDef.Cmd)
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error al ejecutar el comando '%s': %w", cmdDef.Cmd, err)
		}
	}
	return nil
}