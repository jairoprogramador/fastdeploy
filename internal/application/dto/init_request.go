package dto

import "context"

// InitRequest es el DTO para el caso de uso de inicialización del DOM.
type InitRequest struct {
	Ctx              context.Context
	SkipPrompt       bool   // Flag para omitir las preguntas interactivas.
	WorkingDirectory string // El directorio donde se ejecuta 'init'.
}
