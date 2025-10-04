package application

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/newinternal/application/dto"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
)

type ValidateOrderService struct {
}

func NewValidateOrderService() *ValidateOrderService {
	return &ValidateOrderService{}
}

func (s *ValidateOrderService) Validate(
	template *aggregates.DeploymentTemplate, request dto.ValidateOrderRequest) (dto.ValidateOrderResponse, error) {
	environment, err := s.existsEnvironment(template.Environments(), request.Environment)
	if err != nil {
		return dto.ValidateOrderResponse{}, err
	}
	step, err := s.findStep(template.Steps(), request.FinalStep)
	if err != nil {
		return dto.ValidateOrderResponse{}, err
	}

	return dto.ValidateOrderResponse{
		Environment: environment,
		FinalStep:   step.Name(),
	}, nil
}

func (s *ValidateOrderService) existsEnvironment(
	environments []vos.Environment, environmentName string) (vos.Environment, error) {

	var targetEnvironment vos.Environment
	if environmentName == "" {
		if len(environments) == 0 {
			return vos.Environment{},
			fmt.Errorf("no hay ambientes configurados en la plantilla; se debe configurar al menos uno")
		}
		targetEnvironment = environments[0]
	} else {
		env, err := s.findEnvironment(environments, environmentName)
		if err != nil {
			return vos.Environment{}, err
		}
		targetEnvironment = env
	}
	return targetEnvironment, nil
}

func (s *ValidateOrderService) findEnvironment(
	environments []vos.Environment, environmentName string) (vos.Environment, error) {
	for _, env := range environments {
		if env.Name() == environmentName || env.Value() == environmentName {
			return env, nil
		}
	}
	return vos.Environment{},
	fmt.Errorf("el ambiente '%s' no se encontró en la configuración de la plantilla", environmentName)
}

func (s *ValidateOrderService) findStep(
	steps []entities.StepDefinition, stepName string) (entities.StepDefinition, error) {
	for _, step := range steps {
		if step.Name() == stepName || step.Name()[:1] == stepName {
			return step, nil
		}
	}
	return entities.StepDefinition{},
	fmt.Errorf("el paso '%s' no se encontró en la configuración de la plantilla", stepName)
}