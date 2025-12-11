package ports

import "context"

// CommandRunResult encapsula el resultado de ejecutar un comando.
type CommandRunResult struct {
	// Output es la salida combinada de stdout y stderr.
	Output string
	// ExitCode es el código de salida del comando.
	ExitCode int
}

// CommandRunner define la interfaz para ejecutar comandos del sistema.
// Esto permite abstraer la implementación real (ej. os/exec) para facilitar las pruebas.
type CommandRunner interface {
	// Run ejecuta un comando en un directorio de trabajo específico.
	// Devuelve el resultado de la ejecución o un error si el comando no se puede iniciar.
	// Un código de salida distinto de cero no se considera un error de la interfaz,
	// sino que debe ser manejado por el llamador a través de CommandRunResult.
	Run(ctx context.Context, command string, workDir string) (*CommandRunResult, error)
}
