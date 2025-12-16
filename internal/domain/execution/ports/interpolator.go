package ports

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"

type Interpolator interface {
	Interpolate(input string, vars vos.VariableSet) (string, error)
}
