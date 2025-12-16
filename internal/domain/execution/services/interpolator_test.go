package services_test

import (
	"testing"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/services"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterpolator_Interpolate(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		vars           vos.VariableSet
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "Test Basico Exitoso",
			input:          "Hola, ${var.nombre}!",
			vars:           vos.VariableSet{"nombre": "Mundo"},
			expectedOutput: "Hola, Mundo!",
			expectError:    false,
		},
		{
			name:           "Multiples Variables",
			input:          "El valor de ${var.uno} es 1 y el de ${var.dos} es 2.",
			vars:           vos.VariableSet{"uno": "ONE", "dos": "TWO"},
			expectedOutput: "El valor de ONE es 1 y el de TWO es 2.",
			expectError:    false,
		},
		{
			name:           "Variable Faltante",
			input:          "Hola, ${var.nombre}. ¿Cómo estás?",
			vars:           vos.VariableSet{"otro": "valor"},
			expectedOutput: "Hola, . ¿Cómo estás?",
			expectError:    false,
		},
		{
			name:           "Sin Variables en Input",
			input:          "Esta cadena no tiene variables.",
			vars:           vos.VariableSet{"nombre": "Mundo"},
			expectedOutput: "Esta cadena no tiene variables.",
			expectError:    false,
		},
		{
			name:           "Input Vacio",
			input:          "",
			vars:           vos.VariableSet{"nombre": "Mundo"},
			expectedOutput: "",
			expectError:    false,
		},
		{
			name:           "Error de Interpolacion Incompleta",
			input:          "Esto tiene una variable ${malformada}.",
			vars:           vos.VariableSet{},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name:           "Variable al Inicio y al Final",
			input:          "${var.saludo}, te despides con ${var.despedida}",
			vars:           vos.VariableSet{"saludo": "Hola", "despedida": "Adiós"},
			expectedOutput: "Hola, te despides con Adiós",
			expectError:    false,
		},
		{
			name:           "Mapa de Variables Vacio",
			input:          "El valor es ${var.valor}",
			vars:           vos.VariableSet{},
			expectedOutput: "El valor es ",
			expectError:    false,
		},
	}

	interpolator := services.NewInterpolator()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := interpolator.Interpolate(tc.input, tc.vars)

			if tc.expectError {
				require.Error(t, err, "Se esperaba un error pero no se obtuvo")
			} else {
				require.NoError(t, err, "No se esperaba un error pero se obtuvo uno")
				assert.Equal(t, tc.expectedOutput, result, "El resultado de la interpolación no es el esperado")
			}
		})
	}
}
