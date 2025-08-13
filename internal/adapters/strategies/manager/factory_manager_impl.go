package manager

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/java"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/node"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type FactoryManagerImpl struct {
	registry domain.StrategyRegistry
}

func NewFactoryManager() strategies.FactoryManager {
	registry := domain.NewStrategyRegistry()

	registry.Register("java", java.NewJavaFactory())
	registry.Register("node", node.NewNodeFactory())

	return &FactoryManagerImpl{
		registry: registry,
	}
}

func (a *FactoryManagerImpl) GetFactory(projectTechnology string, repositoryPath string) (domain.StrategyFactory, error) {
	factory, err := a.registry.Get(projectTechnology)
	if err != nil {
		return nil, err
	}

	if configurableFactory, ok := factory.(strategies.StrategyFactory); ok {
		configurableFactory.SetRepositoryPath(repositoryPath)
	}

	return factory, nil
}
