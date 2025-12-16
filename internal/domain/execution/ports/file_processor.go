package ports

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"

type FileProcessor interface {
	Process(absPathsFiles []string, vars vos.VariableSet) error
	Restore() error
}
