package state

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/ports"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/state/mapper"
)

type FileOrderRepository struct {
	pathProjectRootFastDeploy string
}

func NewFileOrderRepository(pathProjectRootFastDeploy string) ports.OrderRepository {
	return &FileOrderRepository{pathProjectRootFastDeploy: pathProjectRootFastDeploy}
}

func (r *FileOrderRepository) Save(order *aggregates.Order, nameProject string) error {
	orderDTO := mapper.OrderToDTO(order)

	data, err := yaml.Marshal(orderDTO)
	if err != nil {
		return fmt.Errorf("error al serializar la orden a YAML: %w", err)
	}

	filePath := filepath.Join(r.pathProjectRootFastDeploy, nameProject, "state.yaml")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("error al crear el directorio para el estado de la orden: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error al guardar el archivo de estado de la orden: %w", err)
	}

	return nil
}
