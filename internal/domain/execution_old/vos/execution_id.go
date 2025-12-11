package vos

import (
	"errors"

	"github.com/google/uuid"
)

type ExecutionID struct {
	value uuid.UUID
}

func NewExecutionID() ExecutionID {
	return ExecutionID{value: uuid.New()}
}

func OrderIDFromString(id string) (ExecutionID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return ExecutionID{}, errors.New("el ID de la orden proporcionado no es un UUID válido")
	}
	return ExecutionID{value: parsedUUID}, nil
}

func (id ExecutionID) String() string {
	return id.value.String()
}

func (id ExecutionID) Equals(other ExecutionID) bool {
	return id.value == other.value
}
