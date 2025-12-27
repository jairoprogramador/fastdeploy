package ports

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"

type VarsRepository interface {
	Get(filePath string) (vos.VariableSet, error)
	Save(filePath string, generatedVars vos.VariableSet) error
}
