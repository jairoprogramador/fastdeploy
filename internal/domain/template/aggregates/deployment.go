package aggregates

import (
	"errors"
	"fmt"

	depEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/entities"
	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"
)

type Deployment struct {
	environments []depVos.Environment
	steps        []depEnt.StepDefinition
}

func NewDeployment(
	environments []depVos.Environment,
	steps []depEnt.StepDefinition) (*Deployment, error) {

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

	return &Deployment{
		//source:       source,
		environments: environments,
		steps:        steps,
	}, nil
}

func (dt *Deployment) Environments() []depVos.Environment {
	envsCopy := make([]depVos.Environment, len(dt.environments))
	copy(envsCopy, dt.environments)
	return envsCopy
}

func (dt *Deployment) Steps() []depEnt.StepDefinition {
	stepsCopy := make([]depEnt.StepDefinition, len(dt.steps))
	copy(stepsCopy, dt.steps)
	return stepsCopy
}

func (dt *Deployment) ExistsStep(stepName string) bool {
	for _, step := range dt.steps {
		if step.Name() == stepName || step.Name()[:1] == stepName {
			return true
		}
	}
	return false
}

func (dt *Deployment) StepName(stepName string) string {
	for _, step := range dt.steps {
		if step.Name() == stepName || step.Name()[:1] == stepName {
			return step.Name()
		}
	}
	return stepName
}
