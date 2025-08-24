package manual

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type ManualTest struct {
	strategy.BaseStrategy
}

func NewManualTest(executor service.ExecutorCmd) domain.StepStrategy {
	return &ManualTest{
		BaseStrategy: strategy.BaseStrategy{
			Executor:       executor,
		},
	}
}

func (s *ManualTest) Execute(ctx deployment.Context) error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Java (Spring Boot)")

	if err := s.ExecuteStep(ctx, constant.StepTest, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Pruebas completadas correctamente.")
	return nil
}
