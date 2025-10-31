package vos

import (
	"testing"
)

func TestNewEnvironment(t *testing.T) {

	testCases := []struct {
		testName    string
		name        string
		val         string
		expectError bool
	}{
		{
			testName:    "new environment",
			name:        "sandbox",
			val:         "sand",
			expectError: false,
		},
		{
			testName:    "new environment with empty name",
			name:        "", // Nombre inválido
			val:         "invalid",
			expectError: true,
		},
		{
			testName:    "new environment with empty value",
			name:        "staging",
			val:         "", // Valor inválido
			expectError: true,
		},
	}

	for _, tc := range testCases {
		// t.Run nos permite ejecutar cada caso como un sub-test.
		t.Run(tc.testName, func(t *testing.T) {
			env, err := NewEnvironment(tc.name, tc.val)

			if tc.expectError {
				if err == nil {
					t.Errorf("Se esperaba un error, pero no se obtuvo ninguno")
				}
			} else {
				if err != nil {
					t.Errorf("No se esperaba un error, pero se obtuvo: %v", err)
				}
				if env.Value() != tc.val {
					t.Errorf("Se esperaba el valor '%s', pero se obtuvo '%s'", tc.val, env.Value())
				}
			}
		})
	}
}

func TestEnvironment_Equals(t *testing.T) {
	env1, _ := NewEnvironment("prod", "prod")
	env2, _ := NewEnvironment("prod", "prod")
	env3, _ := NewEnvironment("staging", "stag")

	if !env1.Equals(env2) {
		t.Errorf("Se esperaba que env1 y env2 fueran iguales, pero no lo son")
	}

	if env1.Equals(env3) {
		t.Errorf("Se esperaba que env1 y env3 fueran diferentes, pero son iguales")
	}
}
