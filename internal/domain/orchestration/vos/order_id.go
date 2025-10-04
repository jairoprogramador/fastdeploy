package vos

import (
	"errors"

	"github.com/google/uuid"
)

// OrderID representa el identificador único de una Orden.
// Es un Objeto de Valor que encapsula un UUID para proporcionar seguridad de tipos.
type OrderID struct {
	value uuid.UUID
}

// NewOrderID genera una nueva y única OrderID.
func NewOrderID() OrderID {
	return OrderID{value: uuid.New()}
}

// OrderIDFromString convierte una cadena a una OrderID.
// Devuelve un error si la cadena no es un UUID válido.
func OrderIDFromString(id string) (OrderID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return OrderID{}, errors.New("el ID de la orden proporcionado no es un UUID válido")
	}
	return OrderID{value: parsedUUID}, nil
}

// String devuelve la representación en cadena de la OrderID.
func (id OrderID) String() string {
	return id.value.String()
}

// Equals comprueba la igualdad de valor con otra OrderID.
func (id OrderID) Equals(other OrderID) bool {
	return id.value == other.value
}
