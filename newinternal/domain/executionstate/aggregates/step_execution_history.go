package aggregates

import (
	"errors"
	"sort"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/vos"
)

const maxHistorySize = 5

// StepExecutionHistory es un Agregado Raíz que gestiona la colección de
// recibos de ejecución para un paso específico.
// Su principal responsabilidad es proteger la invariante del tamaño máximo del historial.
type StepExecutionHistory struct {
	stepName string
	receipts []*ExecutionReceipt
}

// NewStepExecutionHistory crea un nuevo historial para un paso.
func NewStepExecutionHistory(stepName string) (*StepExecutionHistory, error) {
	if stepName == "" {
		return nil, errors.New("el nombre del paso no puede estar vacío")
	}
	return &StepExecutionHistory{
		stepName: stepName,
		receipts: make([]*ExecutionReceipt, 0),
	}, nil
}

// RehydrateStepExecutionHistory reconstruye un historial desde un estado persistido.
// Es utilizado por el repositorio para crear el objeto sin la lógica de creación inicial.
func RehydrateStepExecutionHistory(stepName string, receipts []*ExecutionReceipt) *StepExecutionHistory {
	return &StepExecutionHistory{
		stepName: stepName,
		receipts: receipts,
	}
}

// AddReceipt añade un nuevo recibo al historial, manteniendo la invariante del tamaño.
func (h *StepExecutionHistory) AddReceipt(receipt *ExecutionReceipt) {
	if receipt.StepName() != h.stepName {
		return // No añadir recibos que no pertenecen a este historial
	}

	h.receipts = append(h.receipts, receipt)

	// Ordenar por fecha de creación, del más nuevo al más antiguo.
	sort.Slice(h.receipts, func(i, j int) bool {
		return h.receipts[i].CreatedAt().After(h.receipts[j].CreatedAt())
	})

	// Mantener solo los últimos 'maxHistorySize' recibos.
	if len(h.receipts) > maxHistorySize {
		h.receipts = h.receipts[:maxHistorySize]
	}
}

// FindMatch busca el recibo más reciente que coincida con los fingerprints proporcionados.
// Devuelve el recibo encontrado o nil si no hay coincidencia.
func (h *StepExecutionHistory) FindMatch(codeFp, envFp vos.Fingerprint) *ExecutionReceipt {
	// Como la lista está ordenada del más nuevo al más antiguo, la primera coincidencia es la que buscamos.
	for _, receipt := range h.receipts {
		// La lógica de coincidencia debe ser estricta.
		codeMatch := (receipt.CodeFingerprint() == codeFp)
		envMatch := (receipt.EnvironmentFingerprint() == envFp)

		if codeMatch && envMatch {
			return receipt
		}
	}
	return nil
}

// StepName devuelve el nombre del paso de este historial.
func (h *StepExecutionHistory) StepName() string {
	return h.stepName
}

// Receipts devuelve una copia de los recibos en el historial.
func (h *StepExecutionHistory) Receipts() []*ExecutionReceipt {
	receiptsCopy := make([]*ExecutionReceipt, len(h.receipts))
	copy(receiptsCopy, h.receipts)
	return receiptsCopy
}
