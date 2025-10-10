package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/aggregates"
)

type OrderRepository interface {
	Save(order *aggregates.Order, nameProject string) error
}
