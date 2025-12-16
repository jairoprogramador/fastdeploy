package ports

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"

type OutputExtractor interface {
	Extract(commandOutput string, outputs []vos.CommandOutput) (vos.VariableSet, error)
}
