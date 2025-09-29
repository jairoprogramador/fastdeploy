package state

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/state/mapper"
	//"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	//"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/state/dto"
)

// FileOrderRepository implementa la interfaz ports.OrderRepository utilizando el sistema de archivos.
// Persiste el estado de cada orden como un archivo YAML separado.
type FileOrderRepository struct {
	basePath string
}

// NewFileOrderRepository crea una nueva instancia del repositorio de órdenes.
func NewFileOrderRepository(basePath string) (*FileOrderRepository, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("no se pudo crear el directorio base para el estado de las órdenes: %w", err)
	}
	return &FileOrderRepository{basePath: basePath}, nil
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

func (r *FileOrderRepository) getFileStateName(orderId string) string {
	fileStateName := fmt.Sprintf("state%s.yaml", orderId[0:8])
	return fileStateName
}

// FindByID lee un archivo YAML y lo deserializa para reconstruir un agregado Order.
/* func (r *FileOrderRepository) FindByID(_ context.Context, id vos.OrderID, nameProject string) (*aggregates.Order, error) {
	filePath := filepath.Join(r.basePath, nameProject, r.getFileStateName(id.String()))
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no se encontró la orden con ID %s: %w", id.String(), err)
		}
		return nil, fmt.Errorf("error al leer el archivo de estado de la orden: %w", err)
	}

	var orderDTO dto.OrderDTO
	if err := yaml.Unmarshal(data, &orderDTO); err != nil {
		return nil, fmt.Errorf("error al deserializar el estado de la orden desde YAML: %w", err)
	}

	return mapper.OrderToDomain(orderDTO)
} */
