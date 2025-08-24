package service

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/mapper"
	"gopkg.in/yaml.v3"
)

const PROJECT_FILE_NAME = "deploy.yaml"

type FileRepository struct{}

func NewFileRepository() port.Repository {
	return &FileRepository{}
}

func (pr *FileRepository) Load() (entity.Project, error) {
	filePath, err := pr.getFilePath()
	if err != nil {
		return entity.Project{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return entity.Project{}, errors.New("FileNotFoundError: project file does not exist")
		}
		return entity.Project{}, err
	}

	var dto dto.ProjectDto
	err = yaml.Unmarshal(data, &dto)
	if err != nil {
		return entity.Project{}, err
	}
	return mapper.ToDomain(dto)
}

func (pr *FileRepository) Save(project entity.Project) error {
	filePath, err := pr.getFilePath()
	if err != nil {
		return err
	}

	dto, err := mapper.ToDto(project)
	if err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(dto)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, yamlData, 0644)
}

func (pr *FileRepository) Exists() (bool, error) {
	filePath, err := pr.getFilePath()
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

func (pr *FileRepository) getFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, PROJECT_FILE_NAME), nil
}
