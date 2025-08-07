package strategies

import (
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/deploy"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/packet"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/supply"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/test"
)

type JavaFactory struct{}

func (f *JavaFactory) CreateTestStrategy() test.TestStrategy {
	return test.NewJavaTestStrategy()
}

func (f *JavaFactory) CreateSupplyStrategy() supply.SupplyStrategy {
	return supply.NewJavaSupplyStrategy()
}

func (f *JavaFactory) CreatePackageStrategy() packet.PacketStrategy {
	return packet.NewJavaPacketStrategy()
}

func (f *JavaFactory) CreateDeployStrategy() deploy.DeployStrategy {
	return deploy.NewJavaDeployStrategy()
}
