package steps

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type NodeDeploy struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewNodeDeploy(repositoryPath string, executor executor.ExecutorCmd) steps.DeployStrategy {
	return &NodeDeploy{repositoryPath: repositoryPath, executor: executor}
}

func (s *NodeDeploy) ExecuteDeploy(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando deploy para un proyecto Node.js ")
	return nil
}
