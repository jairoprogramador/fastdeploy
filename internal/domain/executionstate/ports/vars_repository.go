package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

type VarsRepository interface {
	Save(vars []vos.Variable) error
	FindAll() ([]vos.Variable, error)
}