package strategies

import "fmt"

type StrategyRegistry interface {
	Register(technology string, factory StrategyFactory)
	Get(technology string) (StrategyFactory, error)
}

type StrategyRegistryImpl struct {
	factories map[string]StrategyFactory
}

func NewStrategyRegistry() StrategyRegistry {
	return &StrategyRegistryImpl{
		factories: make(map[string]StrategyFactory),
	}
}

func (r *StrategyRegistryImpl) Register(technology string, factory StrategyFactory) {
	r.factories[technology] = factory
}

func (r *StrategyRegistryImpl) Get(technology string) (StrategyFactory, error) {
	factory, exists := r.factories[technology]
	if !exists {
		return nil, fmt.Errorf("tecnología de proyecto no soportada: %s", technology)
	}
	return factory, nil
}
