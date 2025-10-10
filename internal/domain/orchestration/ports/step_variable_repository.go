package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

type StepVariableRepository interface {
	Load(stepName string) ([]vos.Variable, error)
}
