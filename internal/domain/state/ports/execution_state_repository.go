package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/state/aggregates"
)

type ExecutionStateRepository interface {
	FindByStepName(stepName string) (*aggregates.ExecutionState, error)
	Save(state *aggregates.ExecutionState) error
}