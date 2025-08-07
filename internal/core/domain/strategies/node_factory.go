package strategies

import (
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/deploy"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/packet"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/supply"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/test"
)

type NodeFactory struct{}

func (f *NodeFactory) CreateTestStrategy() test.TestStrategy {
	return test.NewNodeTestStrategy()
}

func (f *NodeFactory) CreateSupplyStrategy() supply.SupplyStrategy {
	return supply.NewNodeSupplyStrategy()
}

func (f *NodeFactory) CreatePackageStrategy() packet.PacketStrategy {
	return packet.NewNodePacketStrategy()
}

func (f *NodeFactory) CreateDeployStrategy() deploy.DeployStrategy {
	return deploy.NewNodeDeployStrategy()
}
