package strategies

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli/strategies/deploy"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli/strategies/packet"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli/strategies/supply"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli/strategies/test"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type JavaFactory struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewJavaFactory(repositoryPath string) *JavaFactory {
	return &JavaFactory{repositoryPath: repositoryPath, executor: executor.NewCommandExecutor()}
}

func (f *JavaFactory) CreateTestStrategy() steps.TestStrategy {
	return test.NewJavaTestStrategy(f.repositoryPath, f.executor)
}

func (f *JavaFactory) CreateSupplyStrategy() steps.SupplyStrategy {
	return supply.NewJavaSupplyStrategy(f.repositoryPath, f.executor)
}

func (f *JavaFactory) CreatePackageStrategy() steps.PacketStrategy {
	return packet.NewJavaPacketStrategy(f.repositoryPath, f.executor)
}

func (f *JavaFactory) CreateDeployStrategy() steps.DeployStrategy {
	return deploy.NewJavaDeployStrategy(f.repositoryPath, f.executor)
}
