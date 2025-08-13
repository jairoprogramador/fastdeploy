package steps

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type NodePacket struct {
	strategies.BaseStrategy
}

func NewNodePacket(repositoryPath string, executor executor.ExecutorCmd) domain.Strategy {
	return &NodePacket{
		BaseStrategy: strategies.BaseStrategy{
			RepositoryPath: repositoryPath,
			Executor:       executor,
		},
	}
}

func (s *NodePacket) Execute(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Node.js ")
	return nil
}
