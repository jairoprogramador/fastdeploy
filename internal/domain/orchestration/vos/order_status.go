package vos

// OrderStatus representa el estado general de una Order.
type OrderStatus int

const (
	// OrderStatusInProgress indica que la orden se está ejecutando.
	OrderStatusInProgress OrderStatus = iota
	// OrderStatusSuccessful indica que todos los pasos de la orden se completaron exitosamente.
	OrderStatusSuccessful
	// OrderStatusFailed indica que la orden falló porque uno de sus pasos falló.
	OrderStatusFailed
)

// String devuelve la representación en cadena del estado.
func (s OrderStatus) String() string {
	switch s {
	case OrderStatusInProgress:
		return "En Progreso"
	case OrderStatusSuccessful:
		return "Exitoso"
	case OrderStatusFailed:
		return "Fallido"
	default:
		return "Desconocido"
	}
}

func OrderStatusFromString(status string) OrderStatus {
	switch status {
	case "En Progreso":
		return OrderStatusInProgress
	case "Exitoso":
		return OrderStatusSuccessful
	case "Fallido":
		return OrderStatusFailed
	default:
		return OrderStatus(99)
	}
}
