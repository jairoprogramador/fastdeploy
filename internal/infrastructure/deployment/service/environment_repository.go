package service

import (
	"fmt"
	"os"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"os/user"
	"path/filepath"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/mapper"
	"gopkg.in/yaml.v3"
)

const ENVIRONMENT_FILE_NAME = "config.yaml"
const ENVIRONMENT_DIR = "environments"

type EnvironmentRepository struct {
	port.EnvironmentRepository
}

func NewEnvironmentRepository() port.EnvironmentRepository {
	return &EnvironmentRepository{}
}


func (s *EnvironmentRepository) GetEnvironments(repositoryName string) ([]entity.Environment, error) {
	homeDirPath, err := s.getHomeDirPath()
	if err != nil {
		return []entity.Environment{}, err
	}

	filePath := filepath.Join(homeDirPath, repositoryName, ENVIRONMENT_DIR, ENVIRONMENT_FILE_NAME)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []entity.Environment{}, nil
		}
		return []entity.Environment{}, err
	}

	var dtoList []dto.EnvironmentDto
	err = yaml.Unmarshal(data, &dtoList)
	if err != nil {
		return []entity.Environment{}, err
	}

	return mapper.ToDomainList(dtoList)
}

func (s *EnvironmentRepository) getHomeDirPath() (string, error) {
	if fastDeployHome := os.Getenv("FASTDEPLOY_HOME"); fastDeployHome != "" {
		return fastDeployHome, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, constants.FastDeployDir), nil
}