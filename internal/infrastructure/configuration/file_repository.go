package configuration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"fastdeploy/internal/domain/configuration/vos"

	"gopkg.in/yaml.v3"
)

// FileRepository es una implementación de ConfigurationRepository que lee
// la definición de un trabajo desde archivos YAML en el sistema de archivos.
type FileRepository struct {
	// Podríamos definir aquí los nombres de los archivos a buscar,
	// ej. "fastdeploy.yaml", "project.yaml"
}

// NewFileRepository crea una nueva instancia de FileRepository.
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

// Load busca y parsea los archivos de definición de YAML dentro del workspace.
func (r *FileRepository) Load(ctx context.Context, workspacePath string) (*vos.JobDefinition, error) {
	// Asumimos por ahora que toda la definición está en un único archivo "fastdeploy.yaml"
	// Esto podría expandirse para leer múltiples archivos.
	configFilePath := filepath.Join(workspacePath, "fastdeploy.yaml")

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo de configuración '%s': %w", configFilePath, err)
	}

	var jobDefinition vos.JobDefinition
	if err := yaml.Unmarshal(data, &jobDefinition); err != nil {
		return nil, fmt.Errorf("error al parsear el archivo YAML de configuración: %w", err)
	}

	// Aquí podríamos añadir validaciones para asegurar que la definición es correcta.
	if jobDefinition.Project == nil {
		return nil, fmt.Errorf("el archivo de configuración no contiene la sección 'project'")
	}
	if len(jobDefinition.Steps) == 0 {
		return nil, fmt.Errorf("el archivo de configuración no define ningún 'step'")
	}

	// Guardamos la ruta al archivo de instrucciones para el cálculo de fingerprints.
	jobDefinition.PathsToInstructions = make(map[string]string)
	jobDefinition.PathsToInstructions["main"] = configFilePath

	return &jobDefinition, nil
}
