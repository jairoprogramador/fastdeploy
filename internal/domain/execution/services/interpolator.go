package services

import (
	"fmt"
	"os"
	"strings"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

type Interpolator struct{}

func NewInterpolator() *Interpolator {
	return &Interpolator{}
}

func (i *Interpolator) Interpolate(input string, vars vos.VariableSet) (string, error) {
	sanitizedInput := strings.ReplaceAll(input, "${var.", "${var_")

	mapping := func(key string) string {
		originalKey := strings.Replace(key, "var_", "", 1)

		val, exists := vars[originalKey]
		if !exists {
			return ""
		}
		return val
	}

	result := os.Expand(sanitizedInput, mapping)

	if strings.Contains(result, "${") {
		return "", fmt.Errorf("interpolación incompleta, es posible que falten variables. Resultado: %s", result)
	}

	return result, nil
}
