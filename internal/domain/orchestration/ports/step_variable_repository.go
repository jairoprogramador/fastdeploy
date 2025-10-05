package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

type StepVariableRepository interface {
	Load(environment string, stepName string) ([]vos.Variable, error)
}
