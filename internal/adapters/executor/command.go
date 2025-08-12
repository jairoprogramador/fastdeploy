package executor

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/reader"
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
	if _, err := os.Stat(yamlFilePath); os.IsNotExist(err) {
		fmt.Printf("el archivo de comandos YAML no existe en %s\n", yamlFilePath)
		return nil
	}

	commandConfig := reader.CommandConfig{}
	if err := commandConfig.ReadFileYAML(yamlFilePath); err != nil {
		return fmt.Errorf("error al leer el archivo de comandos YAML en %s: %w", yamlFilePath, err)
	}

	for _, cmdDef := range commandConfig.Commands {
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
