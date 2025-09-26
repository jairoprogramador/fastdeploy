package port

import "github.com/jairoprogramador/fastdeploy/internal/domain/variable/values"

type VariablePort interface {
	Load(pathFile string) ([]values.VariableValue, error)
}