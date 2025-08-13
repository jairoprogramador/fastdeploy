package strategies

type StrategyFactory interface {
	CreateTestStrategy() Strategy
	CreateSupplyStrategy() Strategy
	CreatePackageStrategy() Strategy
	CreateDeployStrategy() Strategy
}
