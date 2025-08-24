package ia

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type IATest struct {
	strategy.BaseStrategy
}

func NewIATest(executor service.ExecutorCmd) domain.StepStrategy {
	return &IATest{
		BaseStrategy: strategy.BaseStrategy{
			Executor:       executor,
		},
	}
}

func (s *IATest) Execute(ctx deployment.Context) error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Node.js (ej. npm test)")
	return nil
}
