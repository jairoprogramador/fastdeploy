package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

type StateRepository interface {
	Get(workspacePath string, step vos.Step) (*aggregates.StateTable, error)
	Save(workspacePath string, stateTable *aggregates.StateTable) error
}
