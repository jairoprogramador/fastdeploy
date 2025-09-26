package vos

import "testing"

func TestCommandStatus_String(t *testing.T) {
	testCases := []struct {
		status   CommandStatus
		expected string
	}{
		{CommandStatusPending, "Pendiente"},
		{CommandStatusSuccessful, "Exitoso"},
		{CommandStatusFailed, "Fallido"},
		{CommandStatus(99), "Desconocido"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if tc.status.String() != tc.expected {
				t.Errorf("Se esperaba '%s', pero se obtuvo '%s'", tc.expected, tc.status.String())
			}
		})
	}
}
