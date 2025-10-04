package dom

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/ports"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/dom/dto"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/dom/mapper"
)

// DomYAMLRepository implementa la interfaz ports.DOMRepository.
type DomYAMLRepository struct {
	filePath string
}

// NewDomYAMLRepository crea una instancia del repositorio DOM.
func NewDomYAMLRepository(workingDir string) ports.DOMRepository {
	dirPath := filepath.Join(workingDir, ".fastdeploy")
	return &DomYAMLRepository{
		filePath: filepath.Join(dirPath, "dom.yaml"),
	}
}

// Save serializa y guarda el agregado DOM a dom.yaml.
func (r *DomYAMLRepository) Save(_ context.Context, dom *aggregates.DeploymentObjectModel) error {
	dto := mapper.ToDTO(dom)
	data, err := yaml.Marshal(dto)
	if err != nil {
		return fmt.Errorf("error al serializar dom.yaml: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(r.filePath), 0755); err != nil {
		return fmt.Errorf("error al crear el directorio para dom.yaml: %w", err)
	}
	return os.WriteFile(r.filePath, data, 0644)
}

// Load lee y deserializa el archivo dom.yaml en el agregado DOM.
func (r *DomYAMLRepository) Load(_ context.Context) (*aggregates.DeploymentObjectModel, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var dto dto.DOMDTO
	if err := yaml.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("error al deserializar dom.yaml: %w", err)
	}

	return mapper.ToDomain(dto)
}
