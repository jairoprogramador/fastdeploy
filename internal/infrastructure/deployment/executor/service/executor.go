package service
/*
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/dto"
)

const variablesDirName = "variables"

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
var varRegex = regexp.MustCompile(`\$\{var\.([^}]+)\}`)

type ExecutorCmd interface {
	Execute(yamlFilePath string, context *values.ContextValue) error
}

type CommandExecutor struct{}

func NewCommandExecutor() ExecutorCmd {
	return &CommandExecutor{}
}

func prepareCommand(cmdTemplate string, context *values.ContextValue) (string, error) {
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

func (e *CommandExecutor) Execute(yamlFilePath string, context *values.ContextValue) error {
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
		environment = "local"
	}
	//else {
	//	environment = fmt.Sprintf("%s.yaml", environment)
	//}

	if err := e.processVariablesFromSubDirIfExists(filepath.Dir(yamlFilePath), variablesDirName, environment, context); err != nil {
		return err
	}

	yamlDir := filepath.Dir(yamlFilePath)

	for _, command := range listCmd.Commands {
		if command.NotExecuteLocal && environment == "local" {
			continue
		}
		fmt.Printf("   -> %s: '%s'\n", command.Name, command.Cmd)

		preparedCmd, err := prepareCommand(command.Cmd, context)
		if err != nil {
			return fmt.Errorf("error preparando comando: %w", err)
		}
		var dirExec string
		cmdDir := "."
		dirExec = ""
		if command.Workdir != "" {
			cmdDir = filepath.Join(yamlDir, command.Workdir)

			projectName, err := context.Get(constants.ProjectName)
			if err != nil {
				return fmt.Errorf("error al obtener el nombre del proyecto: %w", err)
			}
			step, err := context.Get(constants.Step)
			if err != nil {
				return fmt.Errorf("error al obtener el step: %w", err)
			}
			destDir, err := e.copyWorkdir(environment, projectName, step, cmdDir)
			if err != nil {
				return fmt.Errorf("error al copiar el directorio de trabajo para el comando '%s': %w", command.Name, err)
			}
			dirExec = destDir
		}

		if command.Templating.Path != nil {
			dir, err := e.processTemplating(cmdDir, command.Templating, context)
			if err != nil {
				return fmt.Errorf("error al procesar las plantillas para el comando '%s': %w", command.Name, err)
			}
			dirExec = dir
		}

		if dirExec != "" {
			cmdDir = dirExec
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

func (e *CommandExecutor) loadAndProcessVariablesFromFile(dirPath string, fileName string, context *values.ContextValue) error {
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

func (e *CommandExecutor) processVariablesFromSubDirIfExists(parentDir, subDir, fileName string, context *values.ContextValue) error {
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

func (e *CommandExecutor) processVariables(variables dto.VariableListDTO, context *values.ContextValue) error {
	for _, variable := range variables {
		preparedValue, err := prepareCommand(variable.Value, context)
		if err != nil {
			return fmt.Errorf("error al preparar el valor para la variable '%s': %w", variable.Name, err)
		}
		context.Set(variable.Name, preparedValue)
	}
	return nil
}

func (e *CommandExecutor) processTemplateFile(filePath string, context *values.ContextValue) error {
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

func (e *CommandExecutor) processTemplating(workdir string, templating dto.TemplatingDTO, context *values.ContextValue) (string, error) {
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

	environment, err := context.Get(constants.Environment)
	if err != nil {
		return "", err
	}

	return e.buildTemplatingPath(environment, projectName, step, fullWorkdir, ""), nil
}

func (e *CommandExecutor) handleTemplateFile(workdir, path, projectName, step string, context *values.ContextValue) error {
	sourceFilePath := filepath.Join(workdir, path)

	environment, err := context.Get(constants.Environment)
	if err != nil {
		return err
	}

	destDir := e.buildTemplatingPath(environment, projectName, step, workdir, path)

	if err := e.copyFileToTemplateDir(sourceFilePath, destDir); err != nil {
		return err
	}

	destFilePath := filepath.Join(destDir, filepath.Base(sourceFilePath))
	//fmt.Printf("   -> Procesando plantilla: %s\n", destFilePath)
	return e.processTemplateFile(destFilePath, context)
}

func (e *CommandExecutor) copyWorkdir(environment, nameProject, step, workdir string) (string, error) {
	destDir := e.buildTemplatingPath(environment, nameProject, step, workdir, "")

	if err := e.copyDirectoryContents(workdir, destDir); err != nil {
		return "", err
	}
	return destDir, nil
}

func (e *CommandExecutor) buildTemplatingPath(environment, nameProject, step, workdir, path string) string {
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

	pathfinal := filepath.Join(homeDirPath, nameProject, environment, step, subDir, pathDir)
	return pathfinal
}

func (e *CommandExecutor) copyDirectoryContents(sourceDir, destDir string) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("no se pudo calcular la ruta relativa para %s: %w", path, err)
		}
		destPath := filepath.Join(destDir, relPath)

		if _, err := os.Stat(destPath); err == nil {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("error inesperado al verificar la ruta de destino %s: %w", destPath, err)
		}
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		sourceFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("no se pudo abrir el archivo de origen %s: %w", path, err)
		}
		defer sourceFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("no se pudo crear el archivo de destino %s: %w", destPath, err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			return fmt.Errorf("no se pudo copiar el contenido de %s a %s: %w", path, destPath, err)
		}

		return nil
	})
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
 */