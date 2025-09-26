package manual

/* import (
	"fmt"

	values "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
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

func (s *ManualSupply) Execute(ctx *values.ContextValue) error {
	environment, err := ctx.Get(constants.Environment)
	if err != nil {
		return err
	}

	if environment != "local" {
		fmt.Println("Ejecutando el comando: SUPPLY")

		if err := s.ExecuteStep(ctx, constant.StepSupply, s.Executor); err != nil {
			return err
		}

		fmt.Println("  [Estrategia] Supply completado correctamente.")
	}else{
		fmt.Println("  [Estrategia] Supply no se ejecutó porque el entorno es local.")
	}

	return nil
}
 */