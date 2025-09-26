package services

import (
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
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

	// ProcessTemplateFile takes a path to a source template file, a destination path,
	// and a map of variables. It reads the source, interpolates the variables,
	// and writes the result to the destination.
	ProcessTemplateFile(srcPath, destPath string, variables map[string]vos.Variable) error

	// ExtractVariable implements the logic previously defined in the VariableExtractor interface.
	// It uses an OutputProbe's regular expression to find and extract a new variable from a text log.
	ExtractVariable(probe deploymentvos.OutputProbe, text string) (variable vos.Variable, match bool, err error)
}
