package application

import (
	"context"
	"fastdeploy/internal/domain/configuration/ports"
	"fastdeploy/internal/domain/configuration/vos"
)

// ConfigurationLoader es un servicio de aplicación que se encarga
// de orquestar la carga y validación de las definiciones de trabajo.
type ConfigurationLoader struct {
	repo ports.ConfigurationRepository
}

// NewConfigurationLoader crea una nueva instancia de ConfigurationLoader.
func NewConfigurationLoader(repo ports.ConfigurationRepository) *ConfigurationLoader {
	return &ConfigurationLoader{repo: repo}
}

// LoadJobDefinition utiliza el repositorio para cargar la definición del trabajo.
func (cl *ConfigurationLoader) LoadJobDefinition(ctx context.Context, workspacePath string) (*vos.JobDefinition, error) {
	// Aquí podríamos añadir lógica extra en el futuro, como validaciones cruzadas,
	// enriquecimiento de datos, etc., antes de devolver la definición.
	return cl.repo.Load(ctx, workspacePath)
}
