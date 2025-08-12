package packet

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodePacketStrategy struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodePacketStrategy(repositoryPath string, executor executor.ExecutorCmd) steps.PacketStrategy {
	return &NodePacketStrategy{repositoryPath: repositoryPath, executor: executor}
}

func (s *NodePacketStrategy) ExecutePacket(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Node.js ")
	return nil
}
