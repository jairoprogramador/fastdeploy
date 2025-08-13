package java

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/java/steps"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type JavaFactory struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewJavaFactory() strategies.StrategyFactory {
	return &JavaFactory{executor: executor.NewCommandExecutor()}
}

func (f *JavaFactory) SetRepositoryPath(repositoryPath string) {
	f.repositoryPath = repositoryPath
}

func (f *JavaFactory) CreateTestStrategy() domain.TestStrategy {
	return steps.NewJavaTest(f.repositoryPath, f.executor)
}

func (f *JavaFactory) CreateSupplyStrategy() domain.SupplyStrategy {
	return steps.NewJavaSupply(f.repositoryPath, f.executor)
}

func (f *JavaFactory) CreatePackageStrategy() domain.PacketStrategy {
	return steps.NewJavaPacket(f.repositoryPath, f.executor)
}

func (f *JavaFactory) CreateDeployStrategy() domain.DeployStrategy {
	return steps.NewJavaDeploy(f.repositoryPath, f.executor)
}
