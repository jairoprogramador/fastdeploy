package manual

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	contextService "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
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

func (s *ManualPacket) Execute(ctx contextService.Context) error {
	fmt.Println("Ejecutando el comando: PACKAGE")

	if err := s.ExecuteStep(ctx, constant.StepPackage, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Package completado correctamente.")
	return nil
}
