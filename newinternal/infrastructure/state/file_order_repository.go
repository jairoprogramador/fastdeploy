package state

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/state/mapper"
)

// FileOrderRepository implementa la interfaz ports.OrderRepository utilizando el sistema de archivos.
// Persiste el estado de cada orden como un archivo YAML separado.
type FileOrderRepository struct {
	basePath string
}

// NewFileOrderRepository crea una nueva instancia del repositorio de órdenes.
func NewFileOrderRepository(basePath string) *FileOrderRepository {
	return &FileOrderRepository{basePath: basePath}
}

// Save serializa el agregado Order a un archivo YAML.
func (r *FileOrderRepository) Save(_ context.Context, order *aggregates.Order, nameProject string) error {
	orderDTO := mapper.OrderToDTO(order)

	data, err := yaml.Marshal(orderDTO)
	if err != nil {
		return fmt.Errorf("error al serializar la orden a YAML: %w", err)
	}

	filePath := filepath.Join(r.basePath, nameProject, "state.yaml")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("error al crear el directorio para el estado de la orden: %w", err)
	}

	//filePath := filepath.Join(r.basePath, nameProject, fmt.Sprintf("%s.yaml", order.ID().String()))
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error al guardar el archivo de estado de la orden: %w", err)
	}

	return nil
}
