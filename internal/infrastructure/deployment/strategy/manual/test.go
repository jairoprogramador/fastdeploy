package manual
/*
import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	values "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type ManualTest struct {
	strategy.BaseStrategy
}

func NewManualTest(executor service.ExecutorCmd) domain.StepStrategy {
	return &ManualTest{
		BaseStrategy: strategy.BaseStrategy{
			Executor:       executor,
		},
	}
}

func (s *ManualTest) Execute(ctx *values.ContextValue) error {
	fmt.Println("Ejecutando el comando: TEST")

	if err := s.ExecuteStep(ctx, constant.StepTest, s.Executor); err != nil {
		return err
	}

	fmt.Println("  [Estrategia] Pruebas completadas correctamente.")
	return nil
}
 */