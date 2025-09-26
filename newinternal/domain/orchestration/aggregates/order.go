package aggregates

import (
	"fmt"

	deployment "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/entities"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/services"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// Order es el Agregado Raíz para el contexto de Orquestación de Ejecución.
// Representa una única invocación del comando de despliegue y protege la consistencia
// de la ejecución de principio a fin.
type Order struct {
	id     vos.OrderID
	status vos.OrderStatus
	targetEnvironment deploymentvos.Environment // Ambiente específico para esta orden.
	stepExecutions []*entities.StepExecution
	variableMap    map[string]vos.Variable
}

// NewOrder es el constructor y guardián del agregado Order.
// Encapsula la lógica de negocio para crear una orden válida y consistente,
// determinando la secuencia de pasos a ejecutar.
func NewOrder(
	id vos.OrderID,
	template *deployment.DeploymentTemplate,
	targetEnvironment deploymentvos.Environment, // Recibe el VO completo.
	finalStepName string,
	skippedStepNames map[string]struct{},
	initialVariables []vos.Variable,
) (*Order, error) {

	definedSteps := template.Steps()
	finalStepIndex := -1
	for i, step := range definedSteps {
		if step.Name() == finalStepName {
			finalStepIndex = i
			break
		}
	}

	if finalStepIndex == -1 {
		return nil, fmt.Errorf("el paso final '%s' no existe en la definición de la plantilla", finalStepName)
	}

	var stepExecutions []*entities.StepExecution
	// Se incluyen todos los pasos desde el principio hasta el paso final.
	for i := 0; i <= finalStepIndex; i++ {
		stepDef := definedSteps[i]
		stepExec, err := entities.NewStepExecution(stepDef)
		if err != nil {
			return nil, fmt.Errorf("error al crear la ejecución para el paso '%s': %w", stepDef.Name(), err)
		}

		if _, shouldSkip := skippedStepNames[stepDef.Name()]; shouldSkip {
			stepExec.Skip()
		}
		stepExecutions = append(stepExecutions, stepExec)
	}

	variableMap := make(map[string]vos.Variable)
	for _, v := range initialVariables {
		variableMap[v.Key()] = v
	}

	// Añadir variables del ambiente al mapa.
	// Usamos un prefijo "env." para evitar colisiones.
	//envNameVar, _ := vos.NewVariable("env.name", targetEnvironment.Name())
	//envValueVar, _ := vos.NewVariable("env.value", targetEnvironment.Value())
	//variableMap[envNameVar.Key()] = envNameVar
	//variableMap[envValueVar.Key()] = envValueVar

	return &Order{
		id:     id,
		status: vos.OrderStatusInProgress,
		//targetEnvironment: targetEnvironment,
		stepExecutions: stepExecutions,
		variableMap:    variableMap,
	}, nil
}

// updateStatus recalcula el estado general de la Orden basándose en sus pasos.
func (o *Order) updateStatus() {
	hasFailed := false
	allCompleted := true
	for _, step := range o.stepExecutions {
		if step.Status() == vos.StepStatusFailed {
			hasFailed = true
			break
		}
		if step.Status() != vos.StepStatusSuccessful && step.Status() != vos.StepStatusSkipped {
			allCompleted = false
		}
	}

	if hasFailed {
		o.status = vos.OrderStatusFailed
	} else if allCompleted {
		o.status = vos.OrderStatusSuccessful
	} else {
		o.status = vos.OrderStatusInProgress
	}
}

// MarkCommandAsCompleted es el método a través del cual el mundo exterior informa
// a la Orden sobre el resultado de la ejecución de un comando.
// El agregado se encarga de orquestar la actualización de estado internamente.
func (o *Order) MarkCommandAsCompleted(
	stepName, commandName, resolvedCmd, log string,
	exitCode int,
	extractor services.VariableResolver,
) error {
	var targetStep *entities.StepExecution
	for _, step := range o.stepExecutions {
		if step.Name() == stepName {
			targetStep = step
			break
		}
	}
	if targetStep == nil {
		return fmt.Errorf("no se encontró el paso '%s' en la orden", stepName)
	}

	// Delegamos la ejecución del comando al StepExecution.
	err := targetStep.CompleteCommand(commandName, resolvedCmd, log, exitCode, extractor)
	if err != nil {
		return fmt.Errorf("error al completar el comando '%s' en el paso '%s': %w", commandName, stepName, err)
	}

	// Tras la ejecución, las nuevas variables se añaden al mapa compartido.
	newVars := targetStep.CollectNewVariables(commandName)
	for _, v := range newVars {
		o.variableMap[v.Key()] = v
	}

	// Finalmente, actualizamos el estado general de la Orden.
	o.updateStatus()

	return nil
}

// ID devuelve el identificador de la Orden.
func (o *Order) ID() vos.OrderID {
	return o.id
}

// Status devuelve el estado actual de la Orden.
func (o *Order) Status() vos.OrderStatus {
	return o.status
}

// RehydrateOrder reconstruye un agregado Order desde un estado persistido.
func RehydrateOrder(id vos.OrderID, status vos.OrderStatus, targetEnv deploymentvos.Environment, steps []*entities.StepExecution, varMap map[string]vos.Variable) *Order {
	return &Order{
		id:                id,
		status:            status,
		targetEnvironment: targetEnv,
		stepExecutions:    steps,
		variableMap:       varMap,
	}
}

// TargetEnvironment devuelve el ambiente objetivo de la orden.
func (o *Order) TargetEnvironment() deploymentvos.Environment {
	return o.targetEnvironment
}

// StepExecutions devuelve las ejecuciones de paso de la orden.
func (o *Order) StepExecutions() []*entities.StepExecution {
	return o.stepExecutions
}

// VariableMap devuelve el mapa de variables compartido de la orden.
func (o *Order) VariableMap() map[string]vos.Variable {
	return o.variableMap
}
