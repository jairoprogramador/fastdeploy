package steps

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodeTest struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodeTest(repositoryPath string, executor executor.ExecutorCmd) steps.TestStrategy {
	return &NodeTest{repositoryPath: repositoryPath, executor: executor}
}

func (s *NodeTest) ExecuteTest(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Node.js (ej. npm test)")
	return nil
}
