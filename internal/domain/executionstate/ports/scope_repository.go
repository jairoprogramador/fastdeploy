package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
)

type ScopeRepository interface {
	SaveEnvironmentStateHistory(history *aggregates.ScopeReceiptHistory, environmentName string, stepName string) error
	SaveCodeStateHistory(history *aggregates.ScopeReceiptHistory) error

	FindEnvironmentStateHistory(environmentName string, stepName string) (*aggregates.ScopeReceiptHistory, error)
	FindCodeStateHistory() (*aggregates.ScopeReceiptHistory, error)
}
