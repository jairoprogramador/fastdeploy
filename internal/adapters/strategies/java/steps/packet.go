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

type JavaPacket struct {
	strategies.BaseStrategy
}

func NewJavaPacket(repositoryPath string, executor executor.ExecutorCmd) domain.Strategy {
	return &JavaPacket{
		BaseStrategy: strategies.BaseStrategy{
			RepositoryPath: repositoryPath,
			Executor:       executor,
		},
	}
}

func (s *JavaPacket) Execute(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Java")

	if err := utils.ExecuteStepFromFile(ctx, s.RepositoryPath, constants.StepPackage, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Package completado correctamente.")
	return nil
}
