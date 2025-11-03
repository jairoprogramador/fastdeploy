package aggregates

import (
	"errors"
	"fmt"

	defEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/entities"
	defVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
)

type Deployment struct {
	environments []defVos.EnvironmentDefinition
	steps        []defEnt.StepDefinition
}

func NewDeployment(
	environments []defVos.EnvironmentDefinition,
	steps []defEnt.StepDefinition) (*Deployment, error) {

	if len(environments) == 0 {
		return nil, errors.New("la plantilla de despliegue debe tener al menos un ambiente")
	}
	if len(steps) == 0 {
		return nil, errors.New("la plantilla de despliegue debe tener al menos un paso")
	}

	envNames := make([]string, len(environments))
	envValues := make([]string, len(environments))
	for i, env := range environments {
		envNames[i] = env.Name()
		envValues[i] = env.Value()
	}
	stepNames := make([]string, len(steps))
	for i, step := range steps {
		stepNames[i] = step.Name()
	}

	if err := validateNoDuplicates(envNames, "nombre de ambiente"); err != nil {
		return nil, err
	}
	if err := validateNoDuplicates(envValues, "valor de ambiente"); err != nil {
		return nil, err
	}
	if err := validateNoDuplicates(stepNames, "nombre de paso"); err != nil {
		return nil, err
	}

	return &Deployment{
		environments: environments,
		steps:        steps,
	}, nil
}

func (dt *Deployment) Environments() []defVos.EnvironmentDefinition {
	envsCopy := make([]defVos.EnvironmentDefinition, len(dt.environments))
	copy(envsCopy, dt.environments)
	return envsCopy
}

func (dt *Deployment) Steps() []defEnt.StepDefinition {
	stepsCopy := make([]defEnt.StepDefinition, len(dt.steps))
	copy(stepsCopy, dt.steps)
	return stepsCopy
}

func (dt *Deployment) ExistsStep(stepName string) bool {
	if _, exists := dt.findStep(stepName); exists {
		return true
	}
	return false
}

func (dt *Deployment) StepName(stepName string) string {
	if step, exists := dt.findStep(stepName); exists {
		return step.Name()
	}
	return stepName
}

func (dt *Deployment) findStep(name string) (defEnt.StepDefinition, bool) {
	for _, step := range dt.steps {
		if step.Name() == name || (len(step.Name()) > 0 && step.Name()[:1] == name) {
			return step, true
		}
	}
	return defEnt.StepDefinition{}, false
}

func validateNoDuplicates(items []string, itemName string) error {
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if _, exists := seen[item]; exists {
			return fmt.Errorf("%s duplicado: %s", itemName, item)
		}
		seen[item] = struct{}{}
	}
	return nil
}