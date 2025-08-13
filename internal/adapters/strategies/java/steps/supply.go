package steps

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/utils"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

type JavaSupply struct {
	strategies.BaseStrategy
}

func NewJavaSupply(repositoryPath string, executor executor.ExecutorCmd) domain.Strategy {
	return &JavaSupply{
		BaseStrategy: strategies.BaseStrategy{
			RepositoryPath: repositoryPath,
			Executor:       executor,
		},
	}
}

func (s *JavaSupply) Execute(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Java (ej. infraestructura)")

	if err := utils.ExecuteStepFromFile(ctx, s.RepositoryPath, constants.StepSupply, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Supply completado correctamente.")
	return nil
}
