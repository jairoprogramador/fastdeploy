package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
)

type StateRepository interface {
	SaveStepStatus(history aggregates.StateSteps) error
	FindStepStatus() (aggregates.StateSteps, error)
}