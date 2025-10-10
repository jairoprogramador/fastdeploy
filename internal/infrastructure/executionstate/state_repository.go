package executionstate

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/ports"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate/mapper"
)

type StateRepository struct {
	pathStateProjectFastDeploy string
}

func NewStateRepository(
	pathStateRootFastDeploy string,
	projectName string,
	repositoryName string,
	environment string) (ports.StateRepository, error) {

	if pathStateRootFastDeploy == "" {
		return nil, fmt.Errorf("path state root is required")
	}
	if projectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if repositoryName == "" {
		return nil, fmt.Errorf("repository name is required")
	}
	if environment == "" {
		return nil, fmt.Errorf("environment is required")
	}

	pathStateProject := filepath.Join(pathStateRootFastDeploy, projectName, repositoryName, environment)

	return &StateRepository{pathStateProjectFastDeploy: pathStateProject}, nil
}

func (r *StateRepository) SaveStepStatus(stateSteps aggregates.StateSteps) error {
	dto := mapper.StateStepsToDTO(stateSteps)
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(dto); err != nil {
		return fmt.Errorf("error al serializar a formato gob: %w", err)
	}

	stateFilePath := r.getPathFileExecution()

	if err := os.MkdirAll(filepath.Dir(stateFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base: %w", err)
	}

	return os.WriteFile(stateFilePath, buffer.Bytes(), 0644)
}

func (r *StateRepository) FindStepStatus() (aggregates.StateSteps, error) {
	filePath := r.getPathFileExecution()

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

func (r *StateRepository) getPathFileExecution() string {
	return filepath.Join(r.pathStateProjectFastDeploy, "execution.state")
}
