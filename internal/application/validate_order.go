package application

import (
	"fmt"

	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"

	depAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/aggregates"
	depEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/entities"
	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"
)

type ValidateOrderService struct {
}

func NewValidateOrderService() *ValidateOrderService {
	return &ValidateOrderService{}
}

func (s *ValidateOrderService) Validate(
	template *depAgg.DeploymentTemplate, request appDto.ValidateOrderRequest) (appDto.ValidateOrderResponse, error) {
	environment, err := s.existsEnvironment(template.Environments(), request.Environment)
	if err != nil {
		return appDto.ValidateOrderResponse{}, err
	}
	step, err := s.findStep(template.Steps(), request.FinalStep)
	if err != nil {
		return appDto.ValidateOrderResponse{}, err
	}

	return appDto.ValidateOrderResponse{
		Environment: environment,
		FinalStep:   step.Name(),
	}, nil
}

func (s *ValidateOrderService) existsEnvironment(
	environments []depVos.Environment, environmentName string) (depVos.Environment, error) {

	var targetEnvironment depVos.Environment
	if environmentName == "" {
		if len(environments) == 0 {
			return depVos.Environment{},
			fmt.Errorf("no hay ambientes configurados en la plantilla; se debe configurar al menos uno")
		}
		targetEnvironment = environments[0]
	} else {
		env, err := s.findEnvironment(environments, environmentName)
		if err != nil {
			return depVos.Environment{}, err
		}
		targetEnvironment = env
	}
	return targetEnvironment, nil
}

func (s *ValidateOrderService) findEnvironment(
	environments []depVos.Environment, environmentName string) (depVos.Environment, error) {
	for _, env := range environments {
		if env.Name() == environmentName || env.Value() == environmentName {
			return env, nil
		}
	}
	return depVos.Environment{},
	fmt.Errorf("el ambiente '%s' no se encontró en la configuración de la plantilla", environmentName)
}

func (s *ValidateOrderService) findStep(
	steps []depEnt.StepDefinition, stepName string) (depEnt.StepDefinition, error) {
	for _, step := range steps {
		if step.Name() == stepName || step.Name()[:1] == stepName {
			return step, nil
		}
	}
	return depEnt.StepDefinition{},
	fmt.Errorf("el paso '%s' no se encontró en la configuración de la plantilla", stepName)
}