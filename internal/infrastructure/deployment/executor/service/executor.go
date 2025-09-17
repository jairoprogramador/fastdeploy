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
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/dto"
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
	listCmd, err := LoadCmdList(yamlFilePath)
	if err != nil {
		return err
	}

	if err := e.processVariablesFromSubDirIfExists(filepath.Dir(yamlFilePath), "variables", "computed.yaml", context); err != nil {
		return err
	}

	environment, err := context.Get(constants.Environment)
	if err != nil {
		return err
	}
	if environment == "" {
		environment = "local.yaml"
	} else {
		environment = fmt.Sprintf("%s.yaml", environment)
	}

	if err := e.processVariablesFromSubDirIfExists(filepath.Dir(yamlFilePath), "variables", environment, context); err != nil {
		return err
	}

	yamlDir := filepath.Dir(yamlFilePath)

	for _, command := range listCmd.Commands {
		fmt.Printf("   -> %s: '%s'\n", command.Name, command.Cmd)

		preparedCmd, err := prepareCommand(command.Cmd, context)
		if err != nil {
			return fmt.Errorf("error preparando comando: %w", err)
		}

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

		if command.Result != "" {
			matches, err := e.getAllSubMatch(command.Result, out.String())
			if err != nil {
				return err
			}

			if len(matches) == 0 {
				return fmt.Errorf("the command response ('%s') is not fullfill the regex: %s", command.Cmd, command.Result)
			}
		}

		for _, output := range command.Variables {
			matches, err := e.getAllSubMatch(output.Regex, out.String())
			if err != nil {
				return err
			}

			if len(matches) == 0 {
				fmt.Printf("no se encontró coincidencia para el output: %s\n", output.Name)
				fmt.Printf("regex: %s\n", output.Regex)
			}
			for _, m := range matches {
				context.Set(output.Name, m[1])
			}
		}
	}
	return nil
}

func (e *CommandExecutor) getAllSubMatch(regexpresion, info string) ([][]string, error) {
	re, err := regexp.Compile(regexpresion)
	if err != nil {
		return nil, fmt.Errorf("regex invalid: %w", err)
	}

	usefulInfo := ansiRegex.ReplaceAllString(info, "")
	return re.FindAllStringSubmatch(usefulInfo, -1), nil
}

func (e *CommandExecutor) loadAndProcessVariablesFromFile(dirPath string, fileName string, context service.Context) error {
	filePath := filepath.Join(dirPath, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	variables, err := LoadVariableList(filePath)
	if err != nil {
		return fmt.Errorf("error al cargar el archivo de variables '%s': %w", filePath, err)
	}

	if err := e.processVariables(variables, context); err != nil {
		return fmt.Errorf("error al procesar las variables de '%s': %w", filePath, err)
	}

	return nil
}

func (e *CommandExecutor) processVariablesFromSubDirIfExists(parentDir, subDir, fileName string, context service.Context) error {
	dirToConsult := filepath.Join(parentDir, subDir)

	fileInfo, err := os.Stat(dirToConsult)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error al verificar el directorio '%s': %w", dirToConsult, err)
	}
	if !fileInfo.IsDir() {
		return nil
	}
	return e.loadAndProcessVariablesFromFile(dirToConsult, fileName, context)
}

func (e *CommandExecutor) processVariables(variables dto.VariableListDTO, context service.Context) error {
	for _, variable := range variables {
		preparedValue, err := prepareCommand(variable.Value, context)
		if err != nil {
			return fmt.Errorf("error al preparar el valor para la variable '%s': %w", variable.Name, err)
		}
		context.Set(variable.Name, preparedValue)
	}
	return nil
}
