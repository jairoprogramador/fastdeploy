package ports

import "context"

// UserInputProvider define el contrato para un adaptador que puede
// obtener entradas interactivas de un usuario.
type UserInputProvider interface {
	// Prompt muestra una pregunta al usuario y devuelve su respuesta.
	// Si el usuario no introduce nada, se devuelve el valor por defecto.
	Prompt(ctx context.Context, question, defaultValue string) (string, error)

	// Confirm muestra una pregunta de sí/no al usuario.
	Confirm(ctx context.Context, question string, defaultValue bool) (bool, error)
}
