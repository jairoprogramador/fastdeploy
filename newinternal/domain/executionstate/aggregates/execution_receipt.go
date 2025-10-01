package aggregates

import (
	"errors"
	"time"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/vos"
	orchestrationvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// ExecutionReceipt es un Agregado Raíz que actúa como un registro inmutable
// de que un paso específico se completó con éxito para un conjunto de fingerprints.
type ExecutionReceipt struct {
	stepName               string
	codeFingerprint        vos.Fingerprint
	environmentFingerprint vos.Fingerprint
	createdAt              time.Time
	orderID                orchestrationvos.OrderID
}

// NewExecutionReceipt crea un nuevo y válido "recibo" de ejecución.
// Se asegura de que se proporcione al menos un fingerprint.
func NewExecutionReceipt(
	stepName string,
	codeFp vos.Fingerprint,
	envFp vos.Fingerprint,
	orderID orchestrationvos.OrderID,
) (*ExecutionReceipt, error) {
	if stepName == "" {
		return nil, errors.New("el nombre del paso no puede estar vacío")
	}

	// Un recibo debe estar asociado con al menos un fingerprint.
	isCodeFpEmpty := (codeFp == vos.Fingerprint{})
	isEnvFpEmpty := (envFp == vos.Fingerprint{})

	if isCodeFpEmpty && isEnvFpEmpty {
		return nil, errors.New("un recibo de ejecución debe tener al menos un code fingerprint o un environment fingerprint")
	}

	return &ExecutionReceipt{
		stepName:               stepName,
		codeFingerprint:        codeFp,
		environmentFingerprint: envFp,
		createdAt:              time.Now().UTC(),
		orderID:                orderID,
	}, nil
}

// StepName devuelve el nombre del paso al que pertenece este recibo.
func (r *ExecutionReceipt) StepName() string {
	return r.stepName
}

// CodeFingerprint devuelve el fingerprint del código asociado. Puede estar vacío.
func (r *ExecutionReceipt) CodeFingerprint() vos.Fingerprint {
	return r.codeFingerprint
}

// EnvironmentFingerprint devuelve el fingerprint del ambiente asociado. Puede estar vacío.
func (r *ExecutionReceipt) EnvironmentFingerprint() vos.Fingerprint {
	return r.environmentFingerprint
}

// CreatedAt devuelve la fecha y hora de creación del recibo.
func (r *ExecutionReceipt) CreatedAt() time.Time {
	return r.createdAt
}

// OrderID devuelve el ID de la orden que generó este recibo.
func (r *ExecutionReceipt) OrderID() orchestrationvos.OrderID {
	return r.orderID
}
