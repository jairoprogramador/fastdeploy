package aggregates

import (
	"errors"
	"fmt"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
)

// DeploymentTemplate es el Agregado Raíz para el contexto de Definición de Despliegue.
// Representa la "receta" completa y consistente de un despliegue, cargada desde una fuente específica.
type DeploymentTemplate struct {
	source       vos.TemplateSource
	environments []vos.Environment
	steps        []entities.StepDefinition
}

// NewDeploymentTemplate es el constructor para el agregado DeploymentTemplate.
// Actúa como el guardián de las invariantes del agregado, asegurando que
// solo se puedan crear instancias consistentes y válidas.
func NewDeploymentTemplate(source vos.TemplateSource, environments []vos.Environment, steps []entities.StepDefinition) (*DeploymentTemplate, error) {
	if len(environments) == 0 {
		return nil, errors.New("la plantilla de despliegue debe tener al menos un ambiente")
	}
	if len(steps) == 0 {
		return nil, errors.New("la plantilla de despliegue debe tener al menos un paso")
	}

	// Validar que no haya nombres de ambiente duplicados
	envNames := make(map[string]struct{})
	for _, env := range environments {
		if _, exists := envNames[env.Name()]; exists {
			return nil, fmt.Errorf("nombre de ambiente duplicado encontrado: %s", env.Name())
		}
		envNames[env.Name()] = struct{}{}
	}

	// Validar que no haya nombres de paso duplicados
	stepNames := make(map[string]struct{})
	for _, step := range steps {
		if _, exists := stepNames[step.Name()]; exists {
			return nil, fmt.Errorf("nombre de paso duplicado encontrado: %s", step.Name())
		}
		stepNames[step.Name()] = struct{}{}
	}

	// Devolvemos puntero porque los agregados suelen tener un ciclo de vida más complejo.
	return &DeploymentTemplate{
		source:       source,
		environments: environments,
		steps:        steps,
	}, nil
}

// Source devuelve la identidad de la plantilla.
func (dt *DeploymentTemplate) Source() vos.TemplateSource {
	return dt.source
}

// Environments devuelve una copia de los ambientes definidos en la plantilla.
func (dt *DeploymentTemplate) Environments() []vos.Environment {
	envsCopy := make([]vos.Environment, len(dt.environments))
	copy(envsCopy, dt.environments)
	return envsCopy
}

// Steps devuelve una copia de los pasos definidos en la plantilla.
func (dt *DeploymentTemplate) Steps() []entities.StepDefinition {
	stepsCopy := make([]entities.StepDefinition, len(dt.steps))
	copy(stepsCopy, dt.steps)
	return stepsCopy
}
