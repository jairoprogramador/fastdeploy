package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
var varRegex = regexp.MustCompile(`\$\{var\.([^}]+)\}`)

type ExecutorCmd interface {
	Execute(yamlFilePath string, context service.Context) error
}

type CommandExecutor struct{}

func NewCommandExecutor() ExecutorCmd {
	return &CommandExecutor{}
}

func prepareCommand(cmdTemplate string, context service.Context) (string, error) {
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

func (e *CommandExecutor) Execute(yamlFilePath string, context service.Context) error {
	listCmd, err := Load(yamlFilePath)
	if err != nil {
		return err
	}

	yamlDir := filepath.Dir(yamlFilePath)

	for _, command := range listCmd.Commands {
		fmt.Printf("   -> %s\n", command.Name)

		preparedCmd, err := prepareCommand(command.Cmd, context)
		if err != nil {
			return fmt.Errorf("error preparando comando: %w", err)
		}	
		fmt.Printf("   -> command: %s\n", preparedCmd)

			
		cmdDir := "."
		if command.Workdir != "" {
			cmdDir = filepath.Join(yamlDir, command.Workdir)
		}

		commandExec := exec.Command("sh", "-c", preparedCmd)

		var out bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &out)

		commandExec.Dir = cmdDir
		commandExec.Stdout = mw
		commandExec.Stderr = mw

		if err := commandExec.Run(); err != nil {
			if command.ContinueOnError {
				fmt.Printf("error: %v\n", err)
				continue
			} else {
				return fmt.Errorf("'%s': %w", preparedCmd, err)
			}
		}

		for _, output := range command.Outputs {
			re, err := regexp.Compile(output.Regex)
			if err != nil {
				return fmt.Errorf("regex inválida: %w", err)
			}

			outputCmd :=  ansiRegex.ReplaceAllString(out.String(), "")

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
