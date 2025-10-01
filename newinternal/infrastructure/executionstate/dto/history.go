package dto

import (
	"time"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/vos"
	orchestrationvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// DTOs para la serialización del historial de ejecución.

type ExecutionReceiptDTO struct {
	StepName               string
	CodeFingerprint        vos.Fingerprint
	EnvironmentFingerprint vos.Fingerprint
	CreatedAt              time.Time
	OrderID                orchestrationvos.OrderID
}

type StepExecutionHistoryDTO struct {
	StepName string
	Receipts []*ExecutionReceiptDTO
}
