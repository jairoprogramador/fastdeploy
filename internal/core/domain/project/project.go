package project

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func Save(projectEntity ProjectEntity) error {
	data, err := yaml.Marshal(projectEntity)
	if err != nil {
		return fmt.Errorf("error al serializar a YAML: %w", err)
	}

	if err := os.WriteFile(constants.ProjectFileName, data, 0644); err != nil {
		return fmt.Errorf("error al guardar el archivo YAML: %w", err)
	}

	return nil
}

func Load() (*ProjectEntity, error) {
	data, err := os.ReadFile(constants.ProjectFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return &ProjectEntity{}, nil
		}
		return nil, fmt.Errorf("error al leer el archivo de configuración: %w", err)
	}

	var projectEntity ProjectEntity
	if err := yaml.Unmarshal(data, &projectEntity); err != nil {
		return nil, fmt.Errorf("error al deserializar la configuración: %w", err)
	}

	return &projectEntity, nil
}

func getProjectName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio de trabajo: %w", err)
	}
	return filepath.Base(dir), nil
}
