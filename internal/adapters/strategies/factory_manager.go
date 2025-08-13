package strategies

import domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"

type FactoryManager interface {
	GetFactory(projectTechnology string, repositoryPath string) (domain.StrategyFactory, error)
}
