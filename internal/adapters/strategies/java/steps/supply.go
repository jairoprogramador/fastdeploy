package steps

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/utils"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type JavaSupply struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewJavaSupply(repositoryPath string, executor executor.ExecutorCmd) steps.SupplyStrategy {
	return &JavaSupply{repositoryPath: repositoryPath, executor: executor}
}

func (s *JavaSupply) ExecuteSupply(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Java (ej. infraestructura)")

	if err := utils.ExecuteStepFromFile(ctx, s.repositoryPath, constants.StepSupply, s.executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Supply completado correctamente.")
	return nil
}
