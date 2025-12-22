package aggregates

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

type CommandExecutor struct {
	runner          ports.CommandRunner
	fileProcessor   ports.FileProcessor
	interpolator    ports.Interpolator
	outputExtractor ports.OutputExtractor
}

func NewCommandExecutor(
	runner ports.CommandRunner,
	fileProcessor ports.FileProcessor,
	interpolator ports.Interpolator,
	outputExtractor ports.OutputExtractor,
) ports.CommandExecutor {
	return &CommandExecutor{
		runner:          runner,
		fileProcessor:   fileProcessor,
		interpolator:    interpolator,
		outputExtractor: outputExtractor,
	}
}

func (ce *CommandExecutor) Execute(
	ctx context.Context,
	command vos.Command,
	currentVars vos.VariableSet,
	workspaceStep string) *vos.ExecutionResult {

	absPathsFiles := make([]string, len(command.TemplateFiles()))
	for i, filePath := range command.TemplateFiles() {
		absPathsFiles[i] = filepath.Join(workspaceStep, command.Workdir(), filePath)
	}

	if err := ce.fileProcessor.Process(absPathsFiles, currentVars); err != nil {
		return &vos.ExecutionResult{Status: vos.Failure, Error: fmt.Errorf("falló al procesar las plantillas: %w", err)}
	}
	defer ce.fileProcessor.Restore()

	interpolatedCmd, err := ce.interpolator.Interpolate(command.Cmd(), currentVars)
	if err != nil {
		return &vos.ExecutionResult{Status: vos.Failure, Error: fmt.Errorf("falló al interpolar el comando: %w", err)}
	}

	execDir := ""
	if command.Workdir() != "" {
		execDir = filepath.Join(workspaceStep, command.Workdir())
	}
	cmdResult, err := ce.runner.Run(ctx, interpolatedCmd, execDir)
	if err != nil {
		return &vos.ExecutionResult{Status: vos.Failure, Error: fmt.Errorf("no se pudo iniciar el comando: %w", err)}
	}

	if cmdResult.ExitCode != 0 {
		return &vos.ExecutionResult{
			Status: vos.Failure,
			Logs:   cmdResult.CombinedOutput(),
			Error:  fmt.Errorf("el comando %s falló con código de salida %d", interpolatedCmd, cmdResult.ExitCode),
		}
	}

	extractedVars, err := ce.outputExtractor.Extract(cmdResult.NormalizedStdout, command.Outputs())
	if err != nil {
		return &vos.ExecutionResult{
			Status: vos.Failure,
			Logs:   cmdResult.CombinedOutput(),
			Error:  fmt.Errorf("falló al extraer las salidas: %w", err),
		}
	}

	return &vos.ExecutionResult{
		Status:     vos.Success,
		Logs:       cmdResult.CombinedOutput(),
		OutputVars: extractedVars,
	}
}
