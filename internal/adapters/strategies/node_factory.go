package strategies

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/deploy"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/packet"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/supply"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/test"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodeFactory struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodeFactory(repositoryPath string) *NodeFactory {
	return &NodeFactory{repositoryPath: repositoryPath, executor: executor.NewCommandExecutor()}
}

func (f *NodeFactory) CreateTestStrategy() steps.TestStrategy {
	return test.NewNodeTestStrategy(f.repositoryPath, f.executor)
}

func (f *NodeFactory) CreateSupplyStrategy() steps.SupplyStrategy {
	return supply.NewNodeSupplyStrategy(f.repositoryPath, f.executor)
}

func (f *NodeFactory) CreatePackageStrategy() steps.PacketStrategy {
	return packet.NewNodePacketStrategy(f.repositoryPath, f.executor)
}

func (f *NodeFactory) CreateDeployStrategy() steps.DeployStrategy {
	return deploy.NewNodeDeployStrategy(f.repositoryPath, f.executor)
}
