package test

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/utils"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies/steps"
)

type JavaTestStrategy struct {
	repositoryPath string
	executor       executor.ExecutorCmd
}

func NewJavaTestStrategy(repositoryPath string, executor executor.ExecutorCmd) steps.TestStrategy {
	return &JavaTestStrategy{repositoryPath: repositoryPath, executor: executor}
}

func (s *JavaTestStrategy) ExecuteTest(ctx context.Context) error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Java (Spring Boot)")

	if err := utils.ExecuteStepFromFile(ctx, s.repositoryPath, constants.StepTest, s.executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Pruebas completadas correctamente.")
	return nil
}
