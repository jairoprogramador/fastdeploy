package ia

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	contextService "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type IASupply struct {
	strategy.BaseStrategy
}

func NewIASupply(executor service.ExecutorCmd) domain.StepStrategy {
	return &IASupply{
		BaseStrategy: strategy.BaseStrategy{
			Executor:       executor,
		},
	}
}

func (s *IASupply) Execute(ctx contextService.Context) error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Node.js (ej. infraestructura)")
	return nil
}
