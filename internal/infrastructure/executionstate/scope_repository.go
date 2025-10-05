package executionstate

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/ports"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate/mapper"
)

// ScopeRepository implementa la interfaz ports.ScopeRepository usando archivos binarios gob.
type ScopeRepository struct {
	pathStateProjectFastDeploy string
}

// NewScopeRepository crea una nueva instancia del repositorio de historiales.
func NewScopeRepository(pathStateFastDeploy string, projectName string) (ports.ScopeRepository, error) {
	if projectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if pathStateFastDeploy == "" {
		return nil, fmt.Errorf("base path is required")
	}

	pathStateProject := filepath.Join(pathStateFastDeploy, projectName)

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

func (r *ScopeRepository) SaveEnvironmentStateHistory(history *aggregates.ScopeReceiptHistory, environmentName string, stepName string) error {
	dto := mapper.ScopeReceiptHistoryToDTO(history)

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(dto); err != nil {
		return fmt.Errorf("error al serializar el historial a formato gob: %w", err)
	}

	stateFilePath := r.getPathFileEnvironment(environmentName, stepName)

	if err := os.MkdirAll(filepath.Dir(stateFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para el estado de ejecución: %w", err)
	}

	return os.WriteFile(stateFilePath, buffer.Bytes(), 0644)
}

// Find deserializa un historial desde un archivo .state o crea uno nuevo si no existe.
func (r *ScopeRepository) FindEnvironmentStateHistory(environmentName string, stepName string) (*aggregates.ScopeReceiptHistory, error) {
	filePath := r.getPathFileEnvironment(environmentName, stepName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Si el historial no existe, creamos uno nuevo y vacío como dicta el contrato.
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

// Find deserializa un historial desde un archivo .state o crea uno nuevo si no existe.
func (r *ScopeRepository) FindCodeStateHistory() (*aggregates.ScopeReceiptHistory, error) {
	filePath := r.getPathFileCode()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Si el historial no existe, creamos uno nuevo y vacío como dicta el contrato.
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
	return filepath.Join(r.pathStateProjectFastDeploy, "code.state")
}

func (r *ScopeRepository) getPathFileEnvironment(environmentName string, stepName string) string {
	return filepath.Join(r.pathStateProjectFastDeploy, "environment", environmentName, fmt.Sprintf("%s.state", stepName))
}
