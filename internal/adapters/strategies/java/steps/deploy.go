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

type JavaDeploy struct {
	strategies.BaseStrategy
}

func NewJavaDeploy(repositoryPath string, executor executor.ExecutorCmd) domain.Strategy {
	return &JavaDeploy{
		BaseStrategy: strategies.BaseStrategy{
			RepositoryPath: repositoryPath,
			Executor:       executor,
		},
	}
}

func (s *JavaDeploy) Execute(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando deploy para un proyecto Java")

	if err := utils.ExecuteStepFromFile(ctx, s.RepositoryPath, constants.StepDeploy, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Deploy completado correctamente.")
	return nil
}
