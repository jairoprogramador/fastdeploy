package ports

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/aggregates"
	//"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// OrderRepository define el contrato para la persistencia del agregado Order.
// Este puerto permite a la capa de aplicación guardar y recuperar el estado de una ejecución
// sin conocer los detalles de la implementación de almacenamiento (e.g., archivos, base de datos).
type OrderRepository interface {
	// Save guarda el estado actual del agregado Order.
	// La implementación decidirá si es una creación o una actualización.
	Save(ctx context.Context, order *aggregates.Order, nameProject string) error

	// FindByID recupera un agregado Order por su identificador único.
	//FindByID(ctx context.Context, id vos.OrderID, nameProject string) (*aggregates.Order, error)
}
