package executionstate

import (
	"fmt"
	"path/filepath"
	"bytes"
	"encoding/gob"
	"os"

	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/ports"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate/mapper"
)

type StateRepository struct {
	pathStateProjectFastDeploy string
}

func NewStateRepository(pathStateFastDeploy string, projectName string) (ports.StateRepository, error) {
	if projectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if pathStateFastDeploy == "" {
		return nil, fmt.Errorf("base path is required")
	}

	pathStateProject := filepath.Join(pathStateFastDeploy, projectName)

	return &StateRepository{pathStateProjectFastDeploy: pathStateProject}, nil
}

func (r *StateRepository) SaveStateSteps(stateSteps aggregates.StateSteps, environmentName string) error {
	dto := mapper.StateStepsToDTO(stateSteps)
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(dto); err != nil {
		return fmt.Errorf("error al serializar a formato gob: %w", err)
	}

	stateFilePath := r.getPathFileEnvironment(environmentName)

	if err := os.MkdirAll(filepath.Dir(stateFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base: %w", err)
	}

	return os.WriteFile(stateFilePath, buffer.Bytes(), 0644)
}

func (r *StateRepository) FindStateSteps(environmentName string) (aggregates.StateSteps, error) {
	filePath := r.getPathFileEnvironment(environmentName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return aggregates.NewStateSteps(), nil
		}
		return aggregates.NewStateSteps(), fmt.Errorf("error al leer el archivo: %w", err)
	}

	var dto map[string]bool
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&dto); err != nil {
		return aggregates.NewStateSteps(), fmt.Errorf("error al deserializar: %w", err)
	}

	return mapper.StateStepsToDomain(dto)
}

func (r *StateRepository) getPathFileEnvironment(environmentName string) string {
	return filepath.Join(r.pathStateProjectFastDeploy, "steps", fmt.Sprintf("%s.state", environmentName))
}