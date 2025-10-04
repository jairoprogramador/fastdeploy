package shell

import (
	"bytes"
	"io"
	"os"
	"context"
	"os/exec"
	"regexp"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/newinternal/application/ports"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// Executor implementa la interfaz ports.CommandExecutor utilizando el paquete os/exec.
// Es un adaptador que traduce las necesidades de la aplicación a llamadas concretas del sistema operativo.
type Executor struct{}

// NewExecutor crea una nueva instancia del Executor.
func NewExecutor() ports.CommandExecutor {
	return &Executor{}
}

func (e *Executor) CreateWorkDir(workdirs ...string) string {
	return filepath.Join(workdirs...)
}

// Execute ejecuta un comando del sistema, captura su salida combinada (stdout y stderr)
// y gestiona el código de salida. Respeta la cancelación del contexto, lo que permite
// detener comandos de larga duración si es necesario.
func (e *Executor) Execute(ctx context.Context, workdir, command string) (log string, exitCode int, err error) {
	// Usamos exec.CommandContext para que la ejecución del comando respete
	// los timeouts o cancelaciones del contexto.
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	if workdir != "" {
		cmd.Dir = workdir
	}

	var out bytes.Buffer
	multiOutput := io.MultiWriter(os.Stdout, &out)

	cmd.Stdout = multiOutput
	cmd.Stderr = multiOutput

	runErr := cmd.Run()
	log = out.String()

	if runErr != nil {
		// Intentamos obtener el código de salida específico del error.
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			// El comando se ejecutó pero devolvió un código de error (e.g., exit 1).
			// Esto no es un error del sistema, sino un fallo esperado del comando.
			return log, exitErr.ExitCode(), nil
		}
		// Si el error es de otro tipo (e.g., el comando no se encontró),
		// devolvemos un código de salida genérico (-1) y el error para que sea manejado.
		return log, -1, runErr
	}

	log = ansiRegex.ReplaceAllString(log, "")
	// El comando se ejecutó exitosamente.
	return log, 0, nil
}
