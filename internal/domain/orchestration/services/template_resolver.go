package services

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/vos"
)

type TemplateResolver interface {
	ResolveTemplate(template string, outputsMap map[string]vos.Output) (string, error)
	ResolvePath(path string, outputsMap map[string]vos.Output) error
	ResolveOutput(output vos.Output, record string) (outputExtracted vos.Output, match bool, err error)
}
