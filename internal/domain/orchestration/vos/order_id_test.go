package vos

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOrderID(t *testing.T) {
	id := NewOrderID()
	if id.value == uuid.Nil {
		t.Errorf("NewOrderID() generó un UUID nulo, lo cual no es válido")
	}
}

func TestOrderIDFromString(t *testing.T) {
	validUUID := uuid.New()

	testCases := []struct {
		testName    string
		input       string
		expectError bool
	}{
		{
			testName:    "Parseo de UUID valido",
			input:       validUUID.String(),
			expectError: false,
		},
		{
			testName:    "Fallo por string invalido",
			input:       "no-es-un-uuid",
			expectError: true,
		},
		{
			testName:    "Fallo por string vacio",
			input:       "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			id, err := OrderIDFromString(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Se esperaba un error, pero no se obtuvo ninguno")
				}
			} else {
				if err != nil {
					t.Errorf("No se esperaba un error, pero se obtuvo: %v", err)
				}
				if id.String() != tc.input {
					t.Errorf("El ID parseado no coincide con el de entrada. Esperado: %s, Obtenido: %s", tc.input, id.String())
				}
			}
		})
	}
}

func TestOrderID_Equals(t *testing.T) {
	uuidVal := uuid.New()
	id1, _ := OrderIDFromString(uuidVal.String())
	id2, _ := OrderIDFromString(uuidVal.String()) // Mismo valor
	id3 := NewOrderID()                           // Valor diferente

	if !id1.Equals(id2) {
		t.Errorf("Se esperaba que id1 y id2 fueran iguales, pero no lo son")
	}

	if id1.Equals(id3) {
		t.Errorf("Se esperaba que id1 y id3 fueran diferentes, pero son iguales")
	}
}
