package steps

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type NodeDeploy struct {
	strategies.BaseStrategy
}

func NewNodeDeploy(repositoryPath string, executor executor.ExecutorCmd) domain.Strategy {
	return &NodeDeploy{
		BaseStrategy: strategies.BaseStrategy{
			RepositoryPath: repositoryPath,
			Executor:       executor,
		},
	}
}

func (s *NodeDeploy) Execute(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando deploy para un proyecto Node.js ")
	return nil
}
