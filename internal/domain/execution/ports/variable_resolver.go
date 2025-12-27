package ports

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"

type VariableResolver interface {
	Resolve(initialVars, varsToResolve vos.VariableSet) (vos.VariableSet, error)
}