package manual

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type ManualDeploy struct {
	strategy.BaseStrategy
}

func NewManualDeploy(executor service.ExecutorCmd) domain.StepStrategy {
	return &ManualDeploy{
		BaseStrategy: strategy.BaseStrategy{
			Executor:       executor,
		},
	}
}

func (s *ManualDeploy) Execute(ctx deployment.Context) error {
	fmt.Println("  [Estrategia] Ejecutando deploy para un proyecto Java")

	if err := s.ExecuteStep(ctx, constant.StepDeploy, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Deploy completado correctamente.")
	return nil
}
