package executionstate

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/ports"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate/mapper"
)

type ScopeRepository struct {
	pathStateProjectFastDeploy string
}

func NewScopeRepository(
	pathStateRootFastDeploy string,
	projectName string,
	repositoryName string,
	environment string) (ports.ScopeRepository, error) {

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

	return &ScopeRepository{pathStateProjectFastDeploy: pathStateProject}, nil
}

func (r *ScopeRepository) SaveCodeStateHistory(history *aggregates.ScopeReceiptHistory) error {
	dto := mapper.ScopeReceiptHistoryToDTO(history)

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(dto); err != nil {
		return fmt.Errorf("error al serializar el historial a formato gob: %w", err)
	}

	stateFilePath := r.getPathFileCode()

	if err := os.MkdirAll(filepath.Dir(stateFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para el estado de codigo: %w", err)
	}

	return os.WriteFile(stateFilePath, buffer.Bytes(), 0644)
}

func (r *ScopeRepository) SaveStepStateHistory(history *aggregates.ScopeReceiptHistory, stepName string) error {
	dto := mapper.ScopeReceiptHistoryToDTO(history)

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(dto); err != nil {
		return fmt.Errorf("error al serializar el historial a formato gob: %w", err)
	}

	stateFilePath := r.getPathFileStep(stepName)

	if err := os.MkdirAll(filepath.Dir(stateFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para el estado de ejecución: %w", err)
	}

	return os.WriteFile(stateFilePath, buffer.Bytes(), 0644)
}

func (r *ScopeRepository) FindStepStateHistory(stepName string) (*aggregates.ScopeReceiptHistory, error) {
	filePath := r.getPathFileStep(stepName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return aggregates.NewScopeReceiptHistory()
		}
		return nil, fmt.Errorf("error al leer el archivo de scope environment: %w", err)
	}

	var dto dto.ScopeReceiptHistoryDTO
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&dto); err != nil {
		return nil, fmt.Errorf("error al deserializar el scope environment: %w", err)
	}

	return mapper.ScopeReceiptHistoryToDomain(dto), nil
}

func (r *ScopeRepository) FindCodeStateHistory() (*aggregates.ScopeReceiptHistory, error) {
	filePath := r.getPathFileCode()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return aggregates.NewScopeReceiptHistory()
		}
		return nil, fmt.Errorf("error al leer el archivo de scope code: %w", err)
	}

	var dto dto.ScopeReceiptHistoryDTO
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&dto); err != nil {
		return nil, fmt.Errorf("error al deserializar el scope code: %w", err)
	}

	return mapper.ScopeReceiptHistoryToDomain(dto), nil
}

func (r *ScopeRepository) getPathFileCode() string {
	return filepath.Join(r.pathStateProjectFastDeploy, "filesCode.state")
}

func (r *ScopeRepository) getPathFileStep(stepName string) string {
	if len(stepName) > 0 {
		stepName = strings.ToUpper(string(stepName[0])) + strings.ToLower(stepName[1:])
	}
	return filepath.Join(r.pathStateProjectFastDeploy, fmt.Sprintf("files%s.state", stepName))
}
