package ia

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	contextService "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type IADeploy struct {
	strategy.BaseStrategy
}

func NewIADeploy(executor service.ExecutorCmd) domain.StepStrategy {
	return &IADeploy{
		BaseStrategy: strategy.BaseStrategy{
			Executor:       executor,
		},
	}
}

func (s *IADeploy) Execute(ctx contextService.Context) error {
	fmt.Println("  [Estrategia] Ejecutando deploy para un proyecto Node.js ")
	return nil
}
