package dom

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	domAgg "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	domPor "github.com/jairoprogramador/fastdeploy/internal/domain/dom/ports"

	iDomDto "github.com/jairoprogramador/fastdeploy/internal/infrastructure/dom/dto"
	iDomMap "github.com/jairoprogramador/fastdeploy/internal/infrastructure/dom/mapper"
)

type DomYAMLRepository struct {
	filePath string
}

func NewDomYAMLRepository(workingDir string) domPor.DomRepository {
	dirPath := filepath.Join(workingDir, ".fastdeploy")
	return &DomYAMLRepository{
		filePath: filepath.Join(dirPath, "dom.yaml"),
	}
}

func (r *DomYAMLRepository) Save(dom *domAgg.DeploymentObjectModel) error {
	dto := iDomMap.DomToDTO(dom)
	data, err := yaml.Marshal(dto)
	if err != nil {
		return fmt.Errorf("error al serializar dom.yaml: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(r.filePath), 0755); err != nil {
		return fmt.Errorf("error al crear el directorio para dom.yaml: %w", err)
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *DomYAMLRepository) Load() (*domAgg.DeploymentObjectModel, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var dto iDomDto.DomDTO
	if err := yaml.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("error al deserializar dom.yaml: %w", err)
	}

	return iDomMap.DomToDomain(dto)
}
