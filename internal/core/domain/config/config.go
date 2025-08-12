package config

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
	"path/filepath"
)

func Save(configEntity ConfigEntity) error {
	filePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error al crear el directorio de configuración: %w", err)
	}

	data, err := yaml.Marshal(configEntity)
	if err != nil {
		return fmt.Errorf("error al serializar la configuración a YAML: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error al guardar la configuración: %w", err)
	}

	return nil
}

func Load() (*ConfigEntity, error) {
	filePath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ConfigEntity{}, nil
		}
		return nil, fmt.Errorf("error al leer el archivo de configuración: %w", err)
	}

	var cfg ConfigEntity
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error al deserializar la configuración: %w", err)
	}

	return &cfg, nil
}

func GetConfigFilePath() (string, error) {
	configDirPath, err := GetConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDirPath, constants.GlobalConfigFileName), nil
}

func GetConfigDirPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, constants.FastDeployDirName), nil
}
