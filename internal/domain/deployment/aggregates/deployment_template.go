package aggregates

import (
	"errors"
	"fmt"

	sharedVos "github.com/jairoprogramador/fastdeploy/internal/domain/shared/vos"

	depEnt "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entities"
	depVos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
)

type DeploymentTemplate struct {
	source       sharedVos.TemplateSource
	environments []depVos.Environment
	steps        []depEnt.StepDefinition
}

func NewDeploymentTemplate(
	source sharedVos.TemplateSource,
	environments []depVos.Environment,
	steps []depEnt.StepDefinition) (*DeploymentTemplate, error) {

	if len(environments) == 0 {
		return nil, errors.New("la plantilla de despliegue debe tener al menos un ambiente")
	}
	if len(steps) == 0 {
		return nil, errors.New("la plantilla de despliegue debe tener al menos un paso")
	}

	envNames := make(map[string]struct{})
	for _, env := range environments {
		if _, exists := envNames[env.Name()]; exists {
			return nil, fmt.Errorf("nombre de ambiente duplicado: %s", env.Name())
		}
		envNames[env.Name()] = struct{}{}
	}

	envValues := make(map[string]struct{})
	for _, env := range environments {
		if _, exists := envValues[env.Value()]; exists {
			return nil, fmt.Errorf("valor de ambiente duplicado: %s", env.Value())
		}
		envValues[env.Value()] = struct{}{}
	}

	stepNames := make(map[string]struct{})
	for _, step := range steps {
		if _, exists := stepNames[step.Name()]; exists {
			return nil, fmt.Errorf("nombre de paso duplicado: %s", step.Name())
		}
		stepNames[step.Name()] = struct{}{}
	}

	return &DeploymentTemplate{
		source:       source,
		environments: environments,
		steps:        steps,
	}, nil
}

/* func (dt *DeploymentTemplate) SearchStep(stepName string) (*entities.StepDefinition, error) {
	for _, step := range dt.steps {
		if step.Name() == stepName {
			return &step, nil
		}
	}
	return nil, fmt.Errorf("no se encontró el paso '%s' en la plantilla", stepName)
} */

func (dt *DeploymentTemplate) Source() sharedVos.TemplateSource {
	return dt.source
}

func (dt *DeploymentTemplate) Environments() []depVos.Environment {
	envsCopy := make([]depVos.Environment, len(dt.environments))
	copy(envsCopy, dt.environments)
	return envsCopy
}

func (dt *DeploymentTemplate) Steps() []depEnt.StepDefinition {
	stepsCopy := make([]depEnt.StepDefinition, len(dt.steps))
	copy(stepsCopy, dt.steps)
	return stepsCopy
}
