package ports

import (
	"context"
	"fastdeploy/internal/domain/configuration/vos"
)

// ConfigurationRepository define la interfaz para cargar la definición de un trabajo.
type ConfigurationRepository interface {
	// Load lee y valida la configuración desde una fuente (ej. un directorio de workspace)
	// y devuelve un JobDefinition completo.
	Load(ctx context.Context, workspacePath string) (*vos.JobDefinition, error)
}
