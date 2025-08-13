package steps

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodePacket struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodePacket(repositoryPath string, executor executor.ExecutorCmd) steps.PacketStrategy {
	return &NodePacket{repositoryPath: repositoryPath, executor: executor}
}

func (s *NodePacket) ExecutePacket(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Node.js ")
	return nil
}
