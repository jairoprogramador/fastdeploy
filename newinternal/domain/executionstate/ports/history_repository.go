package ports

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/aggregates"
)

// HistoryRepository define el contrato para la persistencia del agregado StepExecutionHistory.
// Este puerto permite a la capa de aplicación guardar y recuperar el historial de un paso
// sin conocer los detalles de la implementación de almacenamiento (e.g., archivos gob, base de datos).
type HistoryRepository interface {
	// Save guarda el estado actual del historial de un paso.
	Save(ctx context.Context, history *aggregates.StepExecutionHistory) error

	// Find recupera el historial de un paso por su nombre.
	// Si no se encuentra un historial previo, puede devolver un agregado nuevo y vacío.
	Find(ctx context.Context, stepName string) (*aggregates.StepExecutionHistory, error)
}
