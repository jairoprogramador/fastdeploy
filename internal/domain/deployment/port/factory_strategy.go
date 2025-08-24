package port

import "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"

type FactoryStrategy interface {
	CreateTestStrategy() strategy.StepStrategy
	CreateSupplyStrategy() strategy.StepStrategy
	CreatePackageStrategy() strategy.StepStrategy
	CreateDeployStrategy() strategy.StepStrategy
}
