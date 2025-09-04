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

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
var varRegex = regexp.MustCompile(`\$\{var\.([^}]+)\}`)
type ExecutorCmd interface {
	Execute(yamlFilePath string, context deployment.Context) error
}

type CommandExecutor struct{}

func NewCommandExecutor() ExecutorCmd {
	return &CommandExecutor{}
}

func cleanANSICodes(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func prepareCommand(cmdTemplate string, context deployment.Context) (string, error) {
	result := varRegex.ReplaceAllStringFunc(cmdTemplate, func(match string) string {
		subMatch := varRegex.FindStringSubmatch(match)
		if len(subMatch) >= 1 {
			value, err := context.Get(subMatch[1])
			if err != nil {
				return match
			}
			if value != "" {
				fmt.Printf("   -> Reemplazando ${var.%s} = %s\n", subMatch[1], value)
				return value
			}
		}
		return match
	})
	return result, nil
}

func (e *CommandExecutor) Execute(yamlFilePath string, context deployment.Context) error {
	listCmd, err := Load(yamlFilePath)
	if err != nil {
		return err
	}

	yamlDir := filepath.Dir(yamlFilePath)

	for _, cmdDef := range listCmd.Commands {
		fmt.Printf("   -> %s\n", cmdDef.Name)

		preparedCmd, err := prepareCommand(cmdDef.Cmd, context)
		if err != nil {
			return fmt.Errorf("error preparando comando: %w", err)
		}	
		fmt.Printf("   -> command: %s\n", preparedCmd)

			
		cmdDir := "."
		if cmdDef.Dir != "" {
			cmdDir = filepath.Join(yamlDir, cmdDef.Dir)
		}

		cmd := exec.Command("sh", "-c", preparedCmd)

		var out bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &out)

		cmd.Dir = cmdDir
		cmd.Stdout = mw
		cmd.Stderr = mw

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("'%s': %w", preparedCmd, err)
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
				fmt.Printf("regex: %s\n", re.String())
			}
			for _, m := range matches {
				context.Set(output.Name, m[1])
			}
		}
	}
	return nil
}
