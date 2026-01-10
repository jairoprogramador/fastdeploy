package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/state/aggregates"
)

type StateRepository interface {
	Get(filePath string) (*aggregates.StateTable, error)
	Save(filePath string, stateTable *aggregates.StateTable) error
}
