package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
)

type ExecutorCmd interface {
	Execute(yamlFilePath string, context deployment.Context) error
}

type CommandExecutor struct{}

func NewCommandExecutor() ExecutorCmd {
	return &CommandExecutor{}
}

func cleanANSICodes(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
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

		cmdDir := "."
		if cmdDef.Dir != "" {
			cmdDir = filepath.Join(yamlDir, cmdDef.Dir)
		}

		cmd := exec.Command("sh", "-c", cmdDef.Cmd)

		var out bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &out)

		cmd.Dir = cmdDir
		cmd.Stdout = mw
		cmd.Stderr = mw

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("'%s': %w", cmdDef.Cmd, err)
		}

		for _, output := range cmdDef.Outputs {
			re, err := regexp.Compile(output.Regex)
			if err != nil {
				return fmt.Errorf("regex inválida: %w", err)
			}

			outputCmd := cleanANSICodes(out.String())

			matches := re.FindAllStringSubmatch(outputCmd, -1)
			if len(matches) == 0 {
				fmt.Printf("no se encontró coincidencia para el output: %s\n", output.Name)
			}
			for _, m := range matches {
				context.Set(output.Name, m[1])
			}
		}
	}
	return nil
}
