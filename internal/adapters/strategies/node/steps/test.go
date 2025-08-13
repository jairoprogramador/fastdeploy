package steps

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type NodeTest struct {
	strategies.BaseStrategy
}

func NewNodeTest(repositoryPath string, executor executor.ExecutorCmd) domain.Strategy {
	return &NodeTest{
		BaseStrategy: strategies.BaseStrategy{
			RepositoryPath: repositoryPath,
			Executor:       executor,
		},
	}
}

func (s *NodeTest) Execute(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Node.js (ej. npm test)")
	return nil
}
