package dom

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"

	domAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
	domPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"

	iDomDto "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project/dto"
	iDomMap "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project/mapper"
)

type DomYAMLRepository struct {
	filePath string
}

func NewDomYAMLRepository(workingDir string) domPor.ConfigRepository {
	return &DomYAMLRepository{
		filePath: filepath.Join(workingDir, "fdconfig.yaml"),
	}
}

func (r *DomYAMLRepository) Save(config *domAgg.Config) error {
	dto := iDomMap.ToDTO(config)
	data, err := yaml.Marshal(dto)
	if err != nil {
		return fmt.Errorf("error al serializar fdconfig.yaml: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(r.filePath), 0755); err != nil {
		return fmt.Errorf("error al crear el directorio para fdconfig.yaml: %w", err)
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *DomYAMLRepository) Load() (*domAgg.Config, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var dto iDomDto.FileConfig
	if err := yaml.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("error al deserializar fdconfig.yaml: %w", err)
	}

	return iDomMap.ToDomain(dto)
}
