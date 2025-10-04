package ports

import (
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/aggregates"
)

type StateRepository interface {
	SaveStateSteps(history aggregates.StateSteps, environmentName string) error
	FindStateSteps(environmentName string) (aggregates.StateSteps, error)
}