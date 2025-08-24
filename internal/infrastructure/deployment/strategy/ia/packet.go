package ia

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
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

func (s *IAPacket) Execute(ctx deployment.Context) error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Node.js ")
	return nil
}
