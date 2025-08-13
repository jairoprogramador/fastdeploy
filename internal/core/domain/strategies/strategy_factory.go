package strategies

import (
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type StrategyFactory interface {
	CreateTestStrategy() steps.TestStrategy
	CreateSupplyStrategy() steps.SupplyStrategy
	CreatePackageStrategy() steps.PacketStrategy
	CreateDeployStrategy() steps.DeployStrategy
}
