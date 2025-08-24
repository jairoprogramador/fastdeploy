package factory

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/port"
)

type RegistryStrategy interface {
	Register(technology string, factory port.FactoryStrategy)
	Get(technology string) (port.FactoryStrategy, error)
}

type RegistryStrategyImpl struct {
	factories map[string]port.FactoryStrategy
}

func NewRegistryStrategy() RegistryStrategy {
	factories := make(map[string]port.FactoryStrategy)
    factories[constants.FactoryManual] = NewManualFactory()
	factories[constants.FactoryIA] = NewIAFactory()

	return &RegistryStrategyImpl{
		factories: factories,
	}
}

func (r *RegistryStrategyImpl) Register(technology string, factory port.FactoryStrategy) {
	r.factories[technology] = factory
}

func (r *RegistryStrategyImpl) Get(technology string) (port.FactoryStrategy, error) {
	factory, exists := r.factories[technology]
	if !exists {
		return nil, fmt.Errorf("tecnología de proyecto no soportada: %s", technology)
	}
	return factory, nil
}
