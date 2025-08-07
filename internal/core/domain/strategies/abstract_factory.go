package strategies

import (
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/deploy"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/packet"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/supply"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/test"
)

type StrategyFactory interface {
	CreateTestStrategy() test.TestStrategy
	CreateSupplyStrategy() supply.SupplyStrategy
	CreatePackageStrategy() packet.PacketStrategy
	CreateDeployStrategy() deploy.DeployStrategy
}
