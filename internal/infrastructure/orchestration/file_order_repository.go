package orchestration

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	orcAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/aggregates"
	orcPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/ports"

	iOrcMap "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/orchestration/mapper"
)

type FileOrderRepository struct {
	myProjectPath string
}

func NewFileOrderRepository(
	rootProjectsPath string,
	projectName string,
	repositoryName string) orcPor.OrderRepository {

	myProjectPath := filepath.Join(rootProjectsPath, projectName, repositoryName)
	return &FileOrderRepository{myProjectPath: myProjectPath}
}

func (r *FileOrderRepository) Save(order *orcAgg.Order) error {
	orderDTO := iOrcMap.OrderToDTO(order)

	data, err := yaml.Marshal(orderDTO)
	if err != nil {
		return fmt.Errorf("error al serializar la orden a YAML: %w", err)
	}

	filePath := filepath.Join(r.myProjectPath, "state.yaml")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("error al crear el directorio para el estado de la orden: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error al guardar el archivo de estado de la orden: %w", err)
	}

	return nil
}
