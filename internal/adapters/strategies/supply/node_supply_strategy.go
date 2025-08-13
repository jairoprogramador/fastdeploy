package supply

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodeSupplyStrategy struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodeSupplyStrategy(repositoryPath string, executor executor.ExecutorCmd) steps.SupplyStrategy {
	return &NodeSupplyStrategy{repositoryPath: repositoryPath, executor: executor}
}

func (s *NodeSupplyStrategy) ExecuteSupply(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Node.js (ej. infraestructura)")
	return nil
}
