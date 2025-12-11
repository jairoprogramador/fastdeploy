package services

import (
	"fmt"
	"regexp"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
	sharedVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
)

type OutputExtractor struct{}

func NewOutputExtractor() *OutputExtractor {
	return &OutputExtractor{}
}

func (oe *OutputExtractor) Extract(outputs []*sharedVos.Output, commandOutput string) (vos.VariableSet, error) {
	extractedVars := make(vos.VariableSet)

	for _, output := range outputs {
		re, err := regexp.Compile(output.Probe)
		if err != nil {
			return nil, fmt.Errorf("expresión regular inválida para la salida '%s': %w", output.Name, err)
		}

		matches := re.FindStringSubmatch(commandOutput)

		if len(matches) < 2 {
			return nil, fmt.Errorf("no se encontró la variable de salida '%s' en la salida del comando. Sonda utilizada: %s", output.Name, output.Probe)
		}

		extractedVars[output.Name] = matches[1]
	}

	return extractedVars, nil
}
