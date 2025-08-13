package steps

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodeSupply struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodeSupply(repositoryPath string, executor executor.ExecutorCmd) steps.SupplyStrategy {
	return &NodeSupply{repositoryPath: repositoryPath, executor: executor}
}

func (s *NodeSupply) ExecuteSupply(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Node.js (ej. infraestructura)")
	return nil
}
