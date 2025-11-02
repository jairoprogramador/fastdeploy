package dom

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	proAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
	proPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"

	iProDto "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project/dto"
	iProMap "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project/mapper"
)

type YamlConfigRepository struct {
}

func NewYamlConfigRepository() proPor.ConfigRepository {
	return &YamlConfigRepository{}
}

func (r *YamlConfigRepository) PathFileConfig(pathProject string) string {
	return filepath.Join(pathProject, "fdconfig.yaml")
}

func (r *YamlConfigRepository) Save(config *proAgg.Config, pathProject string) error {
	pathFileConfig := r.PathFileConfig(pathProject)

	dto := iProMap.ToDTO(config)
	data, err := yaml.Marshal(dto)
	if err != nil {
		return fmt.Errorf("error al serializar fdconfig.yaml: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(pathFileConfig), 0755); err != nil {
		return fmt.Errorf("error al crear el directorio para fdconfig.yaml: %w", err)
	}
	return os.WriteFile(pathFileConfig, data, 0644)
}

func (r *YamlConfigRepository) Load(pathProject string) (*proAgg.Config, error) {
	pathFileConfig := r.PathFileConfig(pathProject)

	data, err := os.ReadFile(pathFileConfig)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var dto iProDto.FileConfig
	if err := yaml.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("error al deserializar fdconfig.yaml: %w", err)
	}

	return iProMap.ToDomain(dto)
}
