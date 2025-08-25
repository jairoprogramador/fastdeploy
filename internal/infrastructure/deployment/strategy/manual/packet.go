package manual

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type ManualPacket struct {
	strategy.BaseStrategy
}

func NewManualPacket(executor service.ExecutorCmd) domain.StepStrategy {
	return &ManualPacket{
		BaseStrategy: strategy.BaseStrategy{
			Executor: executor,
		},
	}
}

func (s *ManualPacket) Execute(ctx deployment.Context) error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Java")

	if err := s.ExecuteStep(ctx, constant.StepPackage, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Package completado correctamente.")
	return nil
}
