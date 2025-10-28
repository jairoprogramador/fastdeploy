package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/vos"
)

type FileLoggerRepository struct {
	logsPath string
}

func NewFileLoggerRepository(logsPath string) (ports.LoggerRepository, error) {
	if err := os.MkdirAll(logsPath, 0755); err != nil {
		return nil, fmt.Errorf("could not create logs directory at %s: %w", logsPath, err)
	}
	return &FileLoggerRepository{logsPath: logsPath}, nil
}

func (r *FileLoggerRepository) Save(log *aggregates.Logger) error {
	filePath := filepath.Join(r.logsPath, log.ID().String()+".yaml")

	data, err := yaml.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal execution log to yaml: %w", err)
	}

	return os.WriteFile(filePath, data, 0644)
}

func (r *FileLoggerRepository) FindByID(id vos.LoggerID) (*aggregates.Logger, error) {
	filePath := filepath.Join(r.logsPath, id.String()+".yaml")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read log file: %w", err)
	}

	var log aggregates.Logger
	if err := yaml.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("failed to unmarshal execution log from yaml: %w", err)
	}

	log.RebuildIndex()

	return &log, nil
}
