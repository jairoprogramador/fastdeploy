package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/logger/mapper"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/logger/dto"
)

type FileLoggerRepository struct {
	pathStateProject string
}

func NewFileLoggerRepository(
	pathStateRootFastDeploy string,
	projectName string,
	repositoryName string) (ports.LoggerRepository, error) {

	pathStateProject := filepath.Join(pathStateRootFastDeploy, projectName, repositoryName, "logs")
	if err := os.MkdirAll(pathStateProject, 0755); err != nil {
		return nil, fmt.Errorf("could not create logs directory at %s: %w", pathStateProject, err)
	}
	return &FileLoggerRepository{pathStateProject: pathStateProject}, nil
}

func (r *FileLoggerRepository) Save(log *aggregates.Logger) error {
	filePath := filepath.Join(r.pathStateProject, "logger.yaml")

	loggerDto := mapper.LoggerToDTO(log)
	data, err := yaml.Marshal(loggerDto)
	if err != nil {
		return fmt.Errorf("failed to marshal execution log to yaml: %w", err)
	}

	return os.WriteFile(filePath, data, 0644)
}

func (r *FileLoggerRepository) Find() (aggregates.Logger, error) {
	filePath := filepath.Join(r.pathStateProject, "logger.yaml")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return aggregates.Logger{}, fmt.Errorf("could not read log file: %w", err)
	}

	var loggerDto dto.LoggerDTO
	if err := yaml.Unmarshal(data, &loggerDto); err != nil {
		return aggregates.Logger{}, fmt.Errorf("failed to unmarshal execution log from yaml: %w", err)
	}

	logger, err := mapper.LoggerToDomain(loggerDto)
	if err != nil {
		return aggregates.Logger{}, err
	}

	return *logger, nil
}
