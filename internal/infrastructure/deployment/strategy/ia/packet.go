package ia

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	contextService "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type IAPacket struct {
	strategy.BaseStrategy
}

func NewIAPacket(executor service.ExecutorCmd) domain.StepStrategy {
	return &IAPacket{
		BaseStrategy: strategy.BaseStrategy{
			Executor:       executor,
		},
	}
}

func (s *IAPacket) Execute(ctx contextService.Context) error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Node.js ")
	return nil
}
