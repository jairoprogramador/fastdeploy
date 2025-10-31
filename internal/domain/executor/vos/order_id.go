package vos

import (
	"errors"

	"github.com/google/uuid"
)

type OrderID struct {
	value uuid.UUID
}

func NewOrderID() OrderID {
	return OrderID{value: uuid.New()}
}

func OrderIDFromString(id string) (OrderID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return OrderID{}, errors.New("el ID de la orden proporcionado no es un UUID válido")
	}
	return OrderID{value: parsedUUID}, nil
}

func (id OrderID) String() string {
	return id.value.String()
}

func (id OrderID) Equals(other OrderID) bool {
	return id.value == other.value
}
