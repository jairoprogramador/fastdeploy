package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
)

type ExecutorCmd interface {
	Execute(yamlFilePath string, context deployment.Context) error
}

type CommandExecutor struct{}

func NewCommandExecutor() ExecutorCmd {
	return &CommandExecutor{}
}

func (e *CommandExecutor) Execute(yamlFilePath string, context deployment.Context) error {
	listCmd, err := Load(yamlFilePath)
	if err != nil {
		return err
	}

	yamlDir := filepath.Dir(yamlFilePath)

	for _, cmdDef := range listCmd.Commands {
		fmt.Printf("   -> %s\n", cmdDef.Name)
		fmt.Printf("   -> command: %s\n", cmdDef.Cmd)

		projectDir := "."
		if cmdDef.Dir != "" {
			projectDir = filepath.Join(yamlDir, cmdDef.Dir)
		}

		cmd := exec.Command("sh", "-c", cmdDef.Cmd)
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("'%s': %w", cmdDef.Cmd, err)
		}
	}
	return nil
}