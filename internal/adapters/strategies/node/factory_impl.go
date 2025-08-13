package node

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/node/steps"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodeFactory struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodeFactory() strategies.StrategyFactory {
	return &NodeFactory{executor: executor.NewCommandExecutor()}
}

func (f *NodeFactory) SetRepositoryPath(repositoryPath string) {
	f.repositoryPath = repositoryPath
}

func (f *NodeFactory) CreateTestStrategy() domain.TestStrategy {
	return steps.NewNodeTest(f.repositoryPath, f.executor)
}

func (f *NodeFactory) CreateSupplyStrategy() domain.SupplyStrategy {
	return steps.NewNodeSupply(f.repositoryPath, f.executor)
}

func (f *NodeFactory) CreatePackageStrategy() domain.PacketStrategy {
	return steps.NewNodePacket(f.repositoryPath, f.executor)
}

func (f *NodeFactory) CreateDeployStrategy() domain.DeployStrategy {
	return steps.NewNodeDeploy(f.repositoryPath, f.executor)
}
