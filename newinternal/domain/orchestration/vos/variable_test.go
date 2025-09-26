package vos

import "testing"

func TestNewVariable(t *testing.T) {
	testCases := []struct {
		testName    string
		key         string
		value       string
		expectError bool
	}{
		{
			testName:    "Creacion de variable valida",
			key:         "project_name",
			value:       "fastdeploy",
			expectError: false,
		},
		{
			testName:    "Fallo por valor vacio",
			key:         "optional_var",
			value:       "",
			expectError: true,
		},
		{
			testName:    "Fallo por clave vacia",
			key:         "",
			value:       "some-value",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			variable, err := NewVariable(tc.key, tc.value)

			if tc.expectError {
				if err == nil {
					t.Errorf("Se esperaba un error, pero no se obtuvo ninguno")
				}
			} else {
				if err != nil {
					t.Errorf("No se esperaba un error, pero se obtuvo: %v", err)
				}
				if variable.Key() != tc.key {
					t.Errorf("Se esperaba la clave '%s', pero se obtuvo '%s'", tc.key, variable.Key())
				}
				if variable.Value() != tc.value {
					t.Errorf("Se esperaba el valor '%s', pero se obtuvo '%s'", tc.value, variable.Value())
				}
			}
		})
	}
}

func TestVariable_Equals(t *testing.T) {
	v1, _ := NewVariable("key1", "value1")
	v2, _ := NewVariable("key1", "value1") // Mismo valor
	v3, _ := NewVariable("key2", "value1") // Clave diferente
	v4, _ := NewVariable("key1", "value2") // Valor diferente

	if !v1.Equals(v2) {
		t.Errorf("Se esperaba que v1 y v2 fueran iguales, pero no lo son")
	}

	if v1.Equals(v3) {
		t.Errorf("Se esperaba que v1 y v3 fueran diferentes, pero son iguales")
	}

	if v1.Equals(v4) {
		t.Errorf("Se esperaba que v1 y v4 fueran diferentes, pero son iguales")
	}
}
