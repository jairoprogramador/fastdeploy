package services

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

type ResolverService interface {
	ResolveString(template string, variables map[string]vos.Output) (string, error)
	ResolvePath(path string, variables map[string]vos.Output) error
	ResolveOutput(output vos.Output, record string) (vos.Output, bool, error)
}
