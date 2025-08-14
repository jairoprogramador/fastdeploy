package repositories

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const CONFIG_FILE_NAME = "config.yaml"

type ConfigRepository struct{}

func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{}
}

// Load lee los datos del archivo especificado por CONFIG_FILE_NAME.
// Si el archivo no existe, debe lanzar una excepción específica (FileNotFoundError).
func (cr *ConfigRepository) Load() (map[string]interface{}, error) {
	data, err := os.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("FileNotFoundError: config file does not exist")
		}
		return nil, err
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Save escribe un diccionario de datos en el archivo CONFIG_FILE_NAME.
func (cr *ConfigRepository) Save(data map[string]interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(CONFIG_FILE_NAME, yamlData, 0644)
}
