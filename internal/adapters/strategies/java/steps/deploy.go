package steps

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/utils"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type JavaDeploy struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewJavaDeploy(repositoryPath string, executor executor.ExecutorCmd) steps.DeployStrategy {
	return &JavaDeploy{repositoryPath: repositoryPath, executor: executor}
}

func (s *JavaDeploy) ExecuteDeploy(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando deploy para un proyecto Java")

	if err := utils.ExecuteStepFromFile(ctx, s.repositoryPath, constants.StepDeploy, s.executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Deploy completado correctamente.")
	return nil
}
