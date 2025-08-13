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

type JavaTest struct {
	strategies.BaseStrategy
}

func NewJavaTest(repositoryPath string, executor executor.ExecutorCmd) domain.Strategy {
	return &JavaTest{
		BaseStrategy: strategies.BaseStrategy{
			RepositoryPath: repositoryPath,
			Executor:       executor,
		},
	}
}

func (s *JavaTest) Execute(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Java (Spring Boot)")

	if err := utils.ExecuteStepFromFile(ctx, s.RepositoryPath, constants.StepTest, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Pruebas completadas correctamente.")
	return nil
}
