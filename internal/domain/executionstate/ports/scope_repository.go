package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
)

type ScopeRepository interface {
	SaveStepStateHistory(history *aggregates.ScopeReceiptHistory, stepName string) error
	SaveCodeStateHistory(history *aggregates.ScopeReceiptHistory) error

	FindStepStateHistory(stepName string) (*aggregates.ScopeReceiptHistory, error)
	FindCodeStateHistory() (*aggregates.ScopeReceiptHistory, error)
}
