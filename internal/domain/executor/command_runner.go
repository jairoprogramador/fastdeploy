package executor

import (
	"context"
	"deploy/internal/infrastructure/tools"
	"fmt"
	"strings"
	"sync"
)

// CommandRunner define la interfaz para ejecutar comandos
type CommandRunner interface {
	Run(ctx context.Context, command string) (string, error)
}

type DefaultCommandRunner struct{}

var (
	instanceCommandRunner *DefaultCommandRunner
	onceCommandRunner     sync.Once
)

// GetCommandRunner retorna la instancia única del ejecutor de comandos
func GetCommandRunner() CommandRunner {
	onceCommandRunner.Do(func() {
		instanceCommandRunner = &DefaultCommandRunner{}
	})
	return instanceCommandRunner
}

// Run ejecuta un comando usando el contexto proporcionado
func (r *DefaultCommandRunner) Run(ctx context.Context, command string) (string, error) {
	// Dividir el comando en partes
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("comando vacío")
	}

	// El primer elemento es el comando, el resto son argumentos
	cmd := parts[0]
	args := parts[1:]

	/* if cmd == "docker" && len(args) > 2 {
		if parts[1] == "compose" {
			cmd = "docker compose"
			args = args[1:]
		}
	} */

	// Ejecutar el comando usando tools.ExecuteCommandWithContext
	return tools.ExecuteCommandWithContext(ctx, cmd, args...)
}
