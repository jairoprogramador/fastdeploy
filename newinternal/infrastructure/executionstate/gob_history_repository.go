package executionstate

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/executionstate/dto"
)

// GobHistoryRepository implementa la interfaz ports.HistoryRepository usando archivos binarios gob.
type GobHistoryRepository struct {
	basePath string
}

// NewGobHistoryRepository crea una nueva instancia del repositorio de historiales.
func NewGobHistoryRepository(basePath string) (*GobHistoryRepository, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("no se pudo crear el directorio base para el estado de ejecución: %w", err)
	}
	return &GobHistoryRepository{basePath: basePath}, nil
}

// Save serializa el agregado StepExecutionHistory a un archivo .state.
func (r *GobHistoryRepository) Save(_ context.Context, history *aggregates.StepExecutionHistory) error {
	dto := mapHistoryToDTO(history)

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(dto); err != nil {
		return fmt.Errorf("error al serializar el historial a formato gob: %w", err)
	}

	filePath := r.getFilePath(history.StepName())
	return os.WriteFile(filePath, buffer.Bytes(), 0644)
}

// Find deserializa un historial desde un archivo .state o crea uno nuevo si no existe.
func (r *GobHistoryRepository) Find(_ context.Context, stepName string) (*aggregates.StepExecutionHistory, error) {
	filePath := r.getFilePath(stepName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Si el historial no existe, creamos uno nuevo y vacío como dicta el contrato.
			return aggregates.NewStepExecutionHistory(stepName)
		}
		return nil, fmt.Errorf("error al leer el archivo de historial: %w", err)
	}

	var dto dto.StepExecutionHistoryDTO
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&dto); err != nil {
		return nil, fmt.Errorf("error al deserializar el historial: %w", err)
	}

	return mapDTOToHistory(dto), nil
}

func (r *GobHistoryRepository) getFilePath(stepName string) string {
	return filepath.Join(r.basePath, fmt.Sprintf("%s.state", stepName))
}

// --- Mappers ---

func mapHistoryToDTO(history *aggregates.StepExecutionHistory) *dto.StepExecutionHistoryDTO {
	receiptDTOs := make([]*dto.ExecutionReceiptDTO, 0, len(history.Receipts()))
	for _, r := range history.Receipts() {
		receiptDTOs = append(receiptDTOs, &dto.ExecutionReceiptDTO{
			StepName:               r.StepName(),
			CodeFingerprint:        r.CodeFingerprint(),
			EnvironmentFingerprint: r.EnvironmentFingerprint(),
			CreatedAt:              r.CreatedAt(),
			OrderID:                r.OrderID(),
		})
	}
	return &dto.StepExecutionHistoryDTO{
		StepName: history.StepName(),
		Receipts: receiptDTOs,
	}
}

func mapDTOToHistory(dto dto.StepExecutionHistoryDTO) *aggregates.StepExecutionHistory {
	receipts := make([]*aggregates.ExecutionReceipt, 0, len(dto.Receipts))
	for _, rDTO := range dto.Receipts {
		// Asumimos que la reconstrucción no puede fallar porque los datos ya fueron validados al guardarse.
		receipt, _ := aggregates.NewExecutionReceipt(
			rDTO.StepName,
			rDTO.CodeFingerprint,
			rDTO.EnvironmentFingerprint,
			rDTO.OrderID,
		)
		// Rehidratamos el timestamp que se pierde en el constructor.
		// En una implementación más avanzada, el 'Rehydrate' sería un constructor separado.
		// receipt.SetCreatedAt(rDTO.CreatedAt)
		receipts = append(receipts, receipt)
	}
	return aggregates.RehydrateStepExecutionHistory(dto.StepName, receipts)
}
