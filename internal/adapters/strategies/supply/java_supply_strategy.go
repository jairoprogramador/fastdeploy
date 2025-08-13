package supply

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/utils"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type JavaSupplyStrategy struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewJavaSupplyStrategy(repositoryPath string, executor executor.ExecutorCmd) steps.SupplyStrategy {
	return &JavaSupplyStrategy{repositoryPath: repositoryPath, executor: executor}
}

func (s *JavaSupplyStrategy) ExecuteSupply(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Java (ej. infraestructura)")

	if err := utils.ExecuteStepFromFile(ctx, s.repositoryPath, constants.StepSupply, s.executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Supply completado correctamente.")
	return nil
}
