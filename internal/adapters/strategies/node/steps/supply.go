package steps

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type NodeSupply struct {
	strategies.BaseStrategy
}

func NewNodeSupply(repositoryPath string, executor executor.ExecutorCmd) domain.Strategy {
	return &NodeSupply{
		BaseStrategy: strategies.BaseStrategy{
			RepositoryPath: repositoryPath,
			Executor:       executor,
		},
	}
}

func (s *NodeSupply) Execute(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Node.js (ej. infraestructura)")
	return nil
}
