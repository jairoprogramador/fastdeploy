package factory
/*
import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy/manual"
)

type ManualFactory struct {
	executor service.ExecutorCmd
}

func NewManualFactory() port.FactoryStrategy {
	return &ManualFactory{executor: service.NewCommandExecutor()}
}

func (f *ManualFactory) CreateTestStrategy() strategy.StepStrategy {
	return manual.NewManualTest(f.executor)
}

func (f *ManualFactory) CreateSupplyStrategy() strategy.StepStrategy {
	return manual.NewManualSupply(f.executor)
}

func (f *ManualFactory) CreatePackageStrategy() strategy.StepStrategy {
	return manual.NewManualPacket(f.executor)
}

func (f *ManualFactory) CreateDeployStrategy() strategy.StepStrategy {
	return manual.NewManualDeploy(f.executor)
}
 */