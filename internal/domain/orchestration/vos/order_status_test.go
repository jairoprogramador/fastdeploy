package vos

import "testing"

func TestOrderStatus_String(t *testing.T) {
	testCases := []struct {
		status   OrderStatus
		expected string
	}{
		{OrderStatusInProgress, "En Progreso"},
		{OrderStatusSuccessful, "Exitoso"},
		{OrderStatusFailed, "Fallido"},
		{OrderStatus(99), "Desconocido"}, // Un valor fuera de rango
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if tc.status.String() != tc.expected {
				t.Errorf("Se esperaba '%s', pero se obtuvo '%s'", tc.expected, tc.status.String())
			}
		})
	}
}
