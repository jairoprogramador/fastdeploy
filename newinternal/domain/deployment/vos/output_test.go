package vos

import "testing"

func TestNewOutputProbe(t *testing.T) {
	testCases := []struct {
		testName    string
		name        string
		description string
		probe       string
		expectError bool
	}{
		{
			testName:    "Creacion valida con nombre",
			name:        "var_name",
			description: "Extracts a variable",
			probe:       "name=(.*)",
			expectError: false,
		},
		{
			testName:    "Creacion valida sin nombre (solo validacion)",
			name:        "",
			description: "Checks for success message",
			probe:       "BUILD SUCCESSFUL",
			expectError: false,
		},
		{
			testName:    "Fallo por sonda vacia",
			name:        "some_name",
			description: "Empty probe should fail",
			probe:       "",
			expectError: true,
		},
		{
			testName:    "Fallo por regex invalida",
			name:        "invalid_regex",
			description: "Invalid regex should fail",
			probe:       "[a-z", // Regex incompleta
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			probe, err := NewOutputProbe(tc.name, tc.description, tc.probe)

			if tc.expectError {
				if err == nil {
					t.Errorf("Se esperaba un error, pero no se obtuvo ninguno")
				}
			} else {
				if err != nil {
					t.Errorf("No se esperaba un error, pero se obtuvo: %v", err)
				}
				if probe.Name() != tc.name {
					t.Errorf("Se esperaba el nombre '%s', pero se obtuvo '%s'", tc.name, probe.Name())
				}
				if probe.Description() != tc.description {
					t.Errorf("Se esperaba la descripción '%s', pero se obtuvo '%s'", tc.description, probe.Description())
				}
				if probe.Probe() != tc.probe {
					t.Errorf("Se esperaba la sonda '%s', pero se obtuvo '%s'", tc.probe, probe.Probe())
				}
			}
		})
	}
}
