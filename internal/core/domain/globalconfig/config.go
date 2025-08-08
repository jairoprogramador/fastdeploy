package globalconfig

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
	"path/filepath"
)

type GlobalConfig struct {
	Organization string `yaml:"organization"`
	TeamName     string `yaml:"teamName"`
	Repository   string `yaml:"repository"`
}

func ConfigFilePath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	configDir := filepath.Join(currentUser.HomeDir, ".fastdeploy")
	return filepath.Join(configDir, "config.yaml"), nil
}

func (c *GlobalConfig) Save() error {
	filePath, err := ConfigFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error al crear el directorio de configuración: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error al serializar la configuración a YAML: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error al guardar la configuración: %w", err)
	}

	return nil
}

func Load() (*GlobalConfig, error) {
	filePath, err := ConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &GlobalConfig{}, nil
		}
		return nil, fmt.Errorf("error al leer el archivo de configuración: %w", err)
	}

	var cfg GlobalConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error al deserializar la configuración: %w", err)
	}

	return &cfg, nil
}
