package test

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodeTestStrategy struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodeTestStrategy(repositoryPath string, executor executor.ExecutorCmd) steps.TestStrategy {
	return &NodeTestStrategy{repositoryPath: repositoryPath, executor: executor}
}

func (s *NodeTestStrategy) ExecuteTest(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Node.js (ej. npm test)")
	return nil
}
