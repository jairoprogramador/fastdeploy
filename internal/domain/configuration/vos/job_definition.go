package vos

import (
	shared_vos "fastdeploy/internal/domain/shared/vos"
)

// ProjectDefinition contiene los metadatos de un proyecto.
type ProjectDefinition struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Team        string `yaml:"team"`
	Description string `yaml:"description"`
}

// JobDefinition es el agregado raíz del contexto de configuración.
// Representa un trabajo completo a ejecutar, con toda su definición y metadatos.
type JobDefinition struct {
	Project *ProjectDefinition           `yaml:"project"`
	Steps   []*shared_vos.StepDefinition `yaml:"steps"`

	// PathsToInstructions podría contener las rutas a los archivos de definición
	// para que el orquestador calcule los fingerprints.
	PathsToInstructions map[string]string `yaml:"-"` // Ignorado por YAML
}
