package packet

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/utils"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type JavaPacketStrategy struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewJavaPacketStrategy(repositoryPath string, executor executor.ExecutorCmd) steps.PacketStrategy {
	return &JavaPacketStrategy{repositoryPath: repositoryPath, executor: executor}
}

func (s *JavaPacketStrategy) ExecutePacket(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Java")

	if err := utils.ExecuteStepFromFile(ctx, s.repositoryPath, constants.StepPackage, s.executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Package completado correctamente.")
	return nil
}
