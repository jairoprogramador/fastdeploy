package ports

import "context"

// CommandExecutor define el contrato para un adaptador que puede ejecutar
// comandos del sistema operativo. La capa de aplicación depende de esta interfaz,
// y la capa de infraestructura proporcionará la implementación concreta (e.g., usando os/exec).
type CommandExecutor interface {
	// Execute ejecuta un comando en un directorio de trabajo específico.
	// Devuelve el log combinado (stdout y stderr) y el código de salida del sistema.
	Execute(ctx context.Context, workdir, command string) (log string, exitCode int, err error)
}
