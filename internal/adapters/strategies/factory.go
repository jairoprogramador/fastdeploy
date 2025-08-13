package strategies

import "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"

type StrategyFactory interface {
	strategies.StrategyFactory
	SetRepositoryPath(repositoryPath string)
}
