package manual

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type ManualSupply struct {
	strategy.BaseStrategy
}

func NewManualSupply(executor service.ExecutorCmd) domain.StepStrategy {
	return &ManualSupply{
		BaseStrategy: strategy.BaseStrategy{
			Executor: executor,
		},
	}
}

func (s *ManualSupply) Execute(ctx deployment.Context) error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto")

	if err := s.ExecuteStep(ctx, constant.StepSupply, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Supply completado correctamente.")
	return nil
}
