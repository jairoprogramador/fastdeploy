package vos

type OrderStatus int

const (
	OrderStatusInProgress OrderStatus = iota
	OrderStatusSuccessful
	OrderStatusFailed
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusInProgress:
		return StatusInProgress.String()
	case OrderStatusSuccessful:
		return StatusSuccessful.String()
	case OrderStatusFailed:
		return StatusFailed.String()
	default:
		return StatusUnknown.String()
	}
}

func OrderStatusFromString(status string) OrderStatus {
	switch status {
	case StatusInProgress.String():
		return OrderStatusInProgress
	case StatusSuccessful.String():
		return OrderStatusSuccessful
	case StatusFailed.String():
		return OrderStatusFailed
	default:
		return OrderStatus(99)
	}
}
