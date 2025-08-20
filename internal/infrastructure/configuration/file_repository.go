package configuration

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/ports"
	"gopkg.in/yaml.v3"
)

const CONFIG_FILE_NAME = "config.yaml"

type FileRepository struct{}

func NewFileRepository() ports.Repository {
	return &FileRepository{}
}

func (cr *FileRepository) Load() (entities.Configuration, error) {
	filePath, err := cr.getFilePath()
	if err != nil {
		return entities.Configuration{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return entities.Configuration{}, errors.New("FileNotFoundError: config file does not exist")
		}
		return entities.Configuration{}, err
	}

	var result entities.Configuration
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return entities.Configuration{}, err
	}

	return result, nil
}

func (cr *FileRepository) Save(config entities.Configuration) error {
	filePath, err := cr.getFilePath()
	if err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, yamlData, 0644)
}

func (cr *FileRepository) Exists() (bool, error) {
	filePath, err := cr.getFilePath()
	if err != nil {
		return false, err
	}

	if _, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (cr *FileRepository) getFilePath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}

	directoryPath := filepath.Join(currentUser.HomeDir, constants.FastDeployDir)

	return filepath.Join(directoryPath, CONFIG_FILE_NAME), nil
}