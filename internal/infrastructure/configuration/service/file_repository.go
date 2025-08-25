package service

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/configuration/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/configuration/mapper"
	"gopkg.in/yaml.v3"
)

const CONFIG_FILE_NAME = "config.yaml"

type FileRepository struct{}

func NewFileRepository() port.Repository {
	return &FileRepository{}
}

func (cr *FileRepository) Load() (entity.Configuration, error) {
	filePath, err := cr.getFilePath()
	if err != nil {
		return entity.Configuration{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return entity.Configuration{}, err
		}
		return entity.Configuration{}, err
	}

	var result dto.ConfigDto
	if err = yaml.Unmarshal(data, &result); err != nil {
		return entity.Configuration{}, err
	}

	return mapper.ToDomain(result)
}

func (cr *FileRepository) Save(config entity.Configuration) error {
	filePath, err := cr.getFilePath()
	if err != nil {
		return err
	}

	dto, err := mapper.ToDto(config)
	if err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(dto)
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
		return "", err
	}

	directoryPath := filepath.Join(currentUser.HomeDir, constants.FastDeployDir)

	return filepath.Join(directoryPath, CONFIG_FILE_NAME), nil
}