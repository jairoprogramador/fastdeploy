package manual

/* import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	values "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/strategy"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
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

func (s *ManualPacket) Execute(ctx *values.ContextValue) error {
	environment, err := ctx.Get(constants.Environment)
	if err != nil {
		return err
	}

	if environment != "local" {
		fmt.Println("Ejecutando el comando: PACKAGE")

		if err := s.ExecuteStep(ctx, constant.StepPackage, s.Executor); err != nil {
			return err
		}

		fmt.Println("  [Estrategia] Package completado correctamente.")
	}else{
		fmt.Println("  [Estrategia] Package no se ejecutó porque el entorno es local.")
	}
	return nil
}
 */