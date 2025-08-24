package factory

import (
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy/ia"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/port"

)

type IAFactory struct {
	executor service.ExecutorCmd
}

func NewIAFactory() port.FactoryStrategy {
	return &IAFactory{executor: service.NewCommandExecutor()}
}

func (f *IAFactory) CreateTestStrategy() strategy.StepStrategy {
	return ia.NewIATest(f.executor)
}

func (f *IAFactory) CreateSupplyStrategy() strategy.StepStrategy {
	return ia.NewIASupply(f.executor)
}

func (f *IAFactory) CreatePackageStrategy() strategy.StepStrategy {
	return ia.NewIAPacket(f.executor)
}

func (f *IAFactory) CreateDeployStrategy() strategy.StepStrategy {
	return ia.NewIADeploy(f.executor)
}

