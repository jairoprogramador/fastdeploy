package vos

import "testing"

func TestStepStatus_String(t *testing.T) {
	testCases := []struct {
		status   StepStatus
		expected string
	}{
		{StepStatusPending, "Pendiente"},
		{StepStatusInProgress, "En Progreso"},
		{StepStatusSkipped, "Omitido"},
		{StepStatusSuccessful, "Exitoso"},
		{StepStatusFailed, "Fallido"},
		{StepStatus(99), "Desconocido"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if tc.status.String() != tc.expected {
				t.Errorf("Se esperaba '%s', pero se obtuvo '%s'", tc.expected, tc.status.String())
			}
		})
	}
}
