package services

import (
	deploymentvos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

// VariableResolver defines the contract for a domain service that handles
// all variable-related operations: interpolation in strings, processing of
// template files, and extraction from command outputs.
// This service is stateless and encapsulates complex business logic that doesn't
// belong in an entity or value object.
type VariableResolver interface {
	// Interpolate takes a string that may contain variable placeholders (e.g., ${var.key})
	// and replaces them with their corresponding values from the provided map.
	Interpolate(template string, variables map[string]vos.Variable) (string, error)

	// ProcessTemplateFile lee un archivo de plantilla (o todos los archivos de un directorio),
	// interpola las variables y escribe las modificaciones.
	ProcessTemplate(path string, variables map[string]vos.Variable) error

	// ExtractVariable implements the logic previously defined in the VariableExtractor interface.
	// It uses an OutputProbe's regular expression to find and extract a new variable from a text log.
	ExtractVariable(probe deploymentvos.OutputProbe, text string) (variable vos.Variable, match bool, err error)
}
