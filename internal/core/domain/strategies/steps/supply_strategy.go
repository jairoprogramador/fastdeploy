package steps

import "github.com/jairoprogramador/fastdeploy/internal/core/domain/context"

type SupplyStrategy interface {
	ExecuteSupply(context.Context) error
}
