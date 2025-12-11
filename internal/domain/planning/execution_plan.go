package planning

import "fastdeploy/internal/domain/state/vos"

// StepAction define la acción a realizar para un paso.
type StepAction string

const (
	// ActionExecute indica que el paso debe ser ejecutado.
	ActionExecute StepAction = "EXECUTE"
	// ActionSkipCached indica que el paso debe ser omitido porque su resultado está en caché.
	ActionSkipCached StepAction = "SKIP_CACHED"
	// ActionSkipFailedDep indica que el paso debe ser omitido porque una de sus dependencias ha fallado.
	ActionSkipFailedDep StepAction = "SKIP_FAILED_DEPENDENCY"
)

// PlannedStep representa un paso dentro de un plan de ejecución, con una acción asignada.
type PlannedStep struct {
	// Step es la definición original del paso.
	Step *vos.Step
	// Action es la acción determinada por el orquestador.
	Action StepAction
	// Reason proporciona contexto sobre por qué se asignó una acción de omisión (skip).
	Reason string
}

// ExecutionPlan representa la lista completa de pasos planificados.
type ExecutionPlan struct {
	// Steps es la secuencia ordenada de pasos a procesar.
	Steps []*PlannedStep
}

// NewExecutionPlan crea un nuevo plan de ejecución.
func NewExecutionPlan(steps []*PlannedStep) *ExecutionPlan {
	return &ExecutionPlan{Steps: steps}
}
