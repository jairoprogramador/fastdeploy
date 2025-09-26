package manual

/* import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	values "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
)

type ManualDeploy struct {
	strategy.BaseStrategy
}

func NewManualDeploy(executor service.ExecutorCmd) domain.StepStrategy {
	return &ManualDeploy{
		BaseStrategy: strategy.BaseStrategy{
			Executor:       executor,
		},
	}
}

func (s *ManualDeploy) Execute(ctx *values.ContextValue) error {
	environment, err := ctx.Get(constants.Environment)
	if err != nil {
		return err
	}
	if environment != "local" {
		fmt.Println("Ejecutando el comando: DEPLOY")

		if err := s.ExecuteStep(ctx, constant.StepDeploy, s.Executor); err != nil {
			return err
		}

		fmt.Println("  [Estrategia] Deploy completado correctamente.")
	}else{
		fmt.Println("  [Estrategia] Deploy no se ejecutó porque el entorno es local.")
	}

	return nil
}
 */