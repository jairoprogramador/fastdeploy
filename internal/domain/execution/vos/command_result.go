package vos

import "strings"

type CommandResult struct {
	RawStdout        string // Salida estándar original y sin procesar.
	RawStderr        string // Salida de error original y sin procesar.
	NormalizedStdout string // Salida estándar normalizada (sin ANSI, trim, etc.).
	NormalizedStderr string // Salida de error normalizada.
	ExitCode         int    // Código de salida del comando.
}

// CombinedOutput devuelve la concatenación de RawStdout y RawStderr.
// Es útil para logging o para mostrar la salida completa al usuario.
func (cr *CommandResult) CombinedOutput() string {
	var builder strings.Builder
	builder.WriteString(cr.RawStdout)
	builder.WriteString(cr.RawStderr)
	return builder.String()
}
