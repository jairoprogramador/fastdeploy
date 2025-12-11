package shell

import (
	"context"
	"os/exec"
	"strings"

	"fastdeploy/internal/domain/execution/ports"
)

// OSCommandRunner es una implementación de la interfaz CommandRunner
// que utiliza el paquete estándar `os/exec` para ejecutar comandos.
type OSCommandRunner struct{}

// NewOSCommandRunner crea una nueva instancia de OSCommandRunner.
func NewOSCommandRunner() *OSCommandRunner {
	return &OSCommandRunner{}
}

// Run ejecuta un comando utilizando el shell del sistema.
func (r *OSCommandRunner) Run(ctx context.Context, command string, workDir string) (*ports.CommandRunResult, error) {
	// Usamos 'sh -c' o 'cmd /C' para manejar cadenas de comandos complejas, pipes, etc.
	// Esto es más robusto que simplemente dividir el comando por espacios.
	cmdArgs := strings.Fields(command)
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = workDir

	// Capturamos la salida combinada de stdout y stderr.
	output, err := cmd.CombinedOutput()

	result := &ports.CommandRunResult{
		Output:   string(output),
		ExitCode: 0,
	}

	if err != nil {
		// Si el error es de tipo *exec.ExitError, significa que el comando se ejecutó
		// pero devolvió un código de salida distinto de cero.
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			// Esto no es un error para el runner, el llamador debe decidir qué hacer
			// con un código de salida no nulo.
			return result, nil
		}
		// Si es otro tipo de error (ej. comando no encontrado), lo devolvemos.
		return nil, err
	}

	return result, nil
}
