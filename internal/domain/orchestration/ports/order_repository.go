package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/aggregates"
)

type OrderRepository interface {
	Save(order *aggregates.Order) error
}
