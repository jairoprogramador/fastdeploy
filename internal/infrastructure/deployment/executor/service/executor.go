package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/user"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/dto"
)

const variablesDirName = "variables"

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
				//fmt.Printf("   -> Reemplazando ${var.%s} = %s\n", subMatch[1], value)
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

	if err := e.processVariablesFromSubDirIfExists(filepath.Dir(yamlFilePath), variablesDirName, "computed.yaml", context); err != nil {
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

	if err := e.processVariablesFromSubDirIfExists(filepath.Dir(yamlFilePath), variablesDirName, environment, context); err != nil {
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

		if command.Templating.Path != nil {
			dir, err := e.processTemplating(cmdDir, command.Templating, context)
			if err != nil {
				return fmt.Errorf("error al procesar las plantillas para el comando '%s': %w", command.Name, err)
			}
			cmdDir = dir
		}

		commandExec := exec.Command("sh", "-c", preparedCmd)

		var out bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &out)

		commandExec.Dir = cmdDir
		fmt.Printf("   -> Command exec dir: %s\n", cmdDir)
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

func (e *CommandExecutor) processTemplateFile(filePath string, context service.Context) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo de plantilla %s: %w", filePath, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		originalLine := scanner.Text()
		processedLine, err := prepareCommand(originalLine, context)
		if err != nil {
			fmt.Printf("   -> error procesando línea de plantilla: %v\n", err)
			lines = append(lines, originalLine)
		} else {
			lines = append(lines, processedLine)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error leyendo el archivo de plantilla %s: %w", filePath, err)
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("no se pudo escribir en el archivo de plantilla %s: %w", filePath, err)
	}

	return nil
}

func (e *CommandExecutor) processTemplating(workdir string, templating dto.TemplatingDTO, context service.Context) (string, error) {
	projectName, err := context.Get(constants.ProjectName)
	if err != nil {
		return workdir, fmt.Errorf("no se pudo obtener el nombre del proyecto del contexto: %w", err)
	}

	step, err := context.Get(constants.Step)
	if err != nil {
		return workdir, fmt.Errorf("no se pudo obtener el step del contexto: %w", err)
	}

	fullWorkdir := workdir

	for _, path := range templating.Path {
		sourcePath := filepath.Join(fullWorkdir, path)

		info, err := os.Stat(sourcePath)
		if err != nil {
			return "", fmt.Errorf("no se pudo acceder a la ruta de la plantilla '%s': %w", sourcePath, err)
		}

		if info.IsDir() {
			err := filepath.Walk(sourcePath, func(currentPath string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if fileInfo.IsDir() {
					return nil
				}

				relativePathInWorkdir, _ := filepath.Rel(fullWorkdir, currentPath)
				return e.handleTemplateFile(fullWorkdir, relativePathInWorkdir, projectName, step, context)
			})
			if err != nil {
				return workdir, fmt.Errorf("error al procesar el directorio de plantillas '%s': %w", sourcePath, err)
			}
		} else {
			if err := e.handleTemplateFile(fullWorkdir, path, projectName, step, context); err != nil {
				return workdir, fmt.Errorf("error al procesar el archivo de plantilla '%s': %w", sourcePath, err)
			}
		}
	}

	return e.buildTemplatingPath(projectName, step, fullWorkdir,""), nil
}

func (e *CommandExecutor) handleTemplateFile(workdir, path, projectName, step string, context service.Context) error {
	sourceFilePath := filepath.Join(workdir, path)
	destDir := e.buildTemplatingPath(projectName, step, workdir, path)

	if err := e.copyFileToTemplateDir(sourceFilePath, destDir); err != nil {
		return err
	}

	destFilePath := filepath.Join(destDir, filepath.Base(sourceFilePath))
	//fmt.Printf("   -> Procesando plantilla: %s\n", destFilePath)
	return e.processTemplateFile(destFilePath, context)
}

func (e *CommandExecutor) buildTemplatingPath(nameProject, step, workdir, path string) string {
	searchPattern := string(os.PathSeparator) + step + string(os.PathSeparator)
	parts := strings.SplitN(workdir, searchPattern, 2)

	var subDir string
	if len(parts) == 2 {
		subDir = parts[1]
	} else {
		subDir = ""
	}

	pathDir := filepath.Dir(path)
	homeDirPath, err := e.getHomeDirPath()
	if err != nil {
		return ""
	}

	pathfinal := filepath.Join(homeDirPath, nameProject, step, subDir, pathDir)
	return pathfinal
}

func (e *CommandExecutor) copyFileToTemplateDir(pathFile, pathDir string) error {
	if err := os.MkdirAll(pathDir, os.ModePerm); err != nil {
		return fmt.Errorf("no se pudo crear el directorio de destino %s: %w", pathDir, err)
	}

	sourceFile, err := os.Open(pathFile)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo de origen %s: %w", pathFile, err)
	}
	defer sourceFile.Close()

	destPath := filepath.Join(pathDir, filepath.Base(pathFile))
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo de destino %s: %w", destPath, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("no se pudo copiar el contenido del archivo de %s a %s: %w", pathFile, destPath, err)
	}

	return nil
}

func (pr *CommandExecutor) getHomeDirPath() (string, error) {
	if fastDeployHome := os.Getenv("FASTDEPLOY_HOME"); fastDeployHome != "" {
		return fastDeployHome, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, constants.FastDeployDir), nil
}
