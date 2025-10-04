package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

type VarsRepository interface {
	Save(vars []vos.Variable, environment string) error
	GetStore(environment string) ([]vos.Variable, error)
}