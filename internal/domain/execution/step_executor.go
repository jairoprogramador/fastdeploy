package execution

import (
	"context"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/services"
	execvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
	sharedVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
	"fmt"
	"path/filepath"
)

// StepExecutor es el agregador que orquesta la ejecución de un único paso.
// Encapsula la lógica de procesamiento de plantillas, ejecución de comandos y extracción de salidas.
type StepExecutor struct {
	runner       ports.CommandRunner
	interpolator *services.Interpolator
}

// NewStepExecutor es el constructor para StepExecutor.
func NewStepExecutor(runner ports.CommandRunner, interpolator *services.Interpolator) *StepExecutor {
	return &StepExecutor{
		runner:       runner,
		interpolator: interpolator,
	}
}

// Execute gestiona el ciclo de vida completo de la ejecución de un paso.
func (se *StepExecutor) Execute(ctx context.Context, step *vos.Step, currentVars execvos.VariableSet, workspaceRoot string) (*execvos.ExecutionResult, error) {
	// 1. Procesamiento de plantillas
	templateProcessor := services.NewTemplateProcessor(workspaceRoot, se.interpolator)
	if err := templateProcessor.Process(step.Templates, currentVars); err != nil {
		return &execvos.ExecutionResult{Status: execvos.Failure, Error: fmt.Errorf("falló al procesar las plantillas: %w", err)}, nil
	}
	// Aseguramos la restauración de las plantillas al final de la ejecución
	defer templateProcessor.Restore()

	// 2. Ejecución del comando
	execDir := filepath.Join(workspaceRoot, step.Workdir)

	interpolatedCmd, err := se.interpolator.Interpolate(step.Cmd, currentVars)
	if err != nil {
		return &execvos.ExecutionResult{Status: execvos.Failure, Error: fmt.Errorf("falló al interpolar el comando: %w", err)}, nil
	}

	cmdResult, err := se.runner.Run(ctx, interpolatedCmd, execDir)
	if err != nil {
		// Este error es si el comando no puede iniciarse.
		return &execvos.ExecutionResult{Status: execvos.Failure, Error: fmt.Errorf("no se pudo iniciar el comando: %w", err)}, nil
	}

	if cmdResult.ExitCode != 0 {
		// Este error es si el comando se ejecuta pero falla (exit code != 0).
		return &execvos.ExecutionResult{
			Status: execvos.Failure,
			Logs:   cmdResult.Output,
			Error:  fmt.Errorf("el comando falló con código de salida %d", cmdResult.ExitCode),
		}, nil
	}

	// 3. Extracción y validación de salidas
	outputExtractor := services.NewOutputExtractor()
	extractedVars, err := outputExtractor.Extract(step.Outputs, cmdResult.Output)
	if err != nil {
		// Este error es si una variable de salida no se encuentra.
		return &execvos.ExecutionResult{
			Status: execvos.Failure,
			Logs:   cmdResult.Output,
			Error:  fmt.Errorf("falló al extraer las salidas: %w", err),
		}, nil
	}

	// ¡��xito!
	return &execvos.ExecutionResult{
		Status:     execvos.Success,
		Logs:       cmdResult.Output,
		OutputVars: extractedVars,
	}, nil
}
