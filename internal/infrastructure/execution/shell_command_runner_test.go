package execution

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellCommandRunner_Run(t *testing.T) {
	runner := NewShellCommandRunner()
	ctx := context.Background()

	t.Run("debería capturar stdout correctamente", func(t *testing.T) {
		cmd := `echo "hello world"`
		result, err := runner.Run(ctx, cmd, "")
		require.NoError(t, err)

		assert.Equal(t, 0, result.ExitCode)
		assert.Contains(t, result.RawStdout, "hello world")
		assert.Equal(t, "hello world", result.NormalizedStdout)
		assert.Empty(t, result.RawStderr)
	})

	t.Run("debería capturar stderr correctamente", func(t *testing.T) {
		var cmd string
		if runtime.GOOS == "windows" {
			cmd = `(>&2 echo error message)`
		} else {
			cmd = `echo "error message" >&2`
		}

		result, err := runner.Run(ctx, cmd, "")
		require.NoError(t, err)

		assert.Equal(t, 0, result.ExitCode)
		assert.Empty(t, result.RawStdout)
		assert.Contains(t, result.RawStderr, "error message")
		assert.Equal(t, "error message", result.NormalizedStderr)
	})

	t.Run("debería capturar un código de salida no cero", func(t *testing.T) {
		cmd := "exit 1"
		if runtime.GOOS == "windows" {
			// En Windows, `exit` cierra el `cmd`, necesitamos una forma diferente.
			// `ver` es un comando que no existe y que devuelve 1.
			// O una forma más explícita:
			cmd = "cmd /c exit 1"
		}

		result, err := runner.Run(ctx, cmd, "")
		require.NoError(t, err, "Se espera un resultado, no un error de ejecución")
		assert.Equal(t, 1, result.ExitCode)
	})

	t.Run("debería manejar un comando inexistente con un código de salida no cero", func(t *testing.T) {
		cmd := "uncomandoquenoexiste12345"
		result, err := runner.Run(ctx, cmd, "")
		require.NoError(t, err, "El runner no debería devolver un error, el error está en el ExitCode")

		assert.NotEqual(t, 0, result.ExitCode, "Se esperaba un código de salida distinto de cero")
		if runtime.GOOS != "windows" {
			assert.Equal(t, 127, result.ExitCode, "En Unix, 127 es el código para 'command not found'")
		}

		assert.Contains(t, result.NormalizedStderr, "not found", "Stderr debería contener un mensaje de 'not found'")
	})

	t.Run("debería ejecutar el comando en el workDir especificado", func(t *testing.T) {
		tmpDir := t.TempDir()

		var cmd string
		if runtime.GOOS == "windows" {
			cmd = "cd"
		} else {
			cmd = "pwd"
		}

		result, err := runner.Run(ctx, cmd, tmpDir)
		require.NoError(t, err)

		// Normalizamos la salida para la comparación, ya que pwd puede tener saltos de línea.
		// `strings.TrimSpace` es suficiente aquí.
		output := strings.TrimSpace(result.RawStdout)

		// `os.Stat` nos da la información del directorio real, que puede ser un enlace simbólico.
		// `filepath.EvalSymlinks` resuelve la ruta real para una comparación fiable.
		expectedDir, err := os.Readlink(tmpDir)
		if err != nil {
			expectedDir = tmpDir
		}

		assert.Contains(t, output, expectedDir, "La salida de pwd/cd debería contener el workDir")
	})
}
