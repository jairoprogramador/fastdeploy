package port

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
)

type LineProcessor func(line string, context *values.ContextValue) string

type TemplatePort interface {
	Process(pathTemplate string, processor LineProcessor, context *values.ContextValue) error
}
