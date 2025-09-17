package manual

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	contextService "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
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

func (s *ManualSupply) Execute(ctx contextService.Context) error {
	fmt.Println("Ejecutando el comando: SUPPLY")

	if err := s.ExecuteStep(ctx, constant.StepSupply, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Supply completado correctamente.")
	return nil
}
