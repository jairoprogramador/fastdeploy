package application
/* import (
	"context"
	"fmt"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/entities"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/aggregates"
	execvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/planning"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/services"
)

// OrchestratorInput encapsula todos los datos necesarios para iniciar una ejecución.
type OrchestratorInput struct {
	Project   *ProjectInfo
	Workspace *WorkspaceInfo
	Steps     []*entities.StepDefinition
}

// ProjectInfo y WorkspaceInfo podrían moverse a VOs de aplicación.
type ProjectInfo struct{ ID, Name, Team string }
type WorkspaceInfo struct{ Path string }

// ExecutionOrchestrator orquesta el flujo completo de ejecución de pasos.
type ExecutionOrchestrator struct {
	stateManager      *services.StateManager
	versionCalculator *versioning.VersionCalculator
	stepExecutor      *aggregates.StepExecutor
}

// NewExecutionOrchestrator es el constructor.
func NewExecutionOrchestrator(
	stateManager *services.StateManager,
	versionCalc *versioning.VersionCalculator,
	stepExec *aggregates.StepExecutor,
) *ExecutionOrchestrator {
	return &ExecutionOrchestrator{
		stateManager:      stateManager,
		versionCalculator: versionCalc,
		stepExecutor:      stepExec,
	}
}

// Run ejecuta el flujo principal de orquestación.
func (o *ExecutionOrchestrator) Run(ctx context.Context, input *OrchestratorInput) (*planning.ExecutionPlan, error) {
	globalVars := o.initializeGlobalVariables(input)
	plan := o.createInitialPlan(input.Steps)
	var firstError error

	for i, plannedStep := range plan.Steps {
		// Si una dependencia falló, este paso ya está marcado como omitido.
		if plannedStep.Action == planning.ActionSkipFailedDep {
			continue
		}

		// 1. Decidir si el paso debe ejecutarse (lógica de estado)
		decision, err := o.stateManager.ShouldExecute(ctx, plannedStep.Step, input.Workspace.Path)
		if err != nil {
			firstError = o.handleStepError(plan, i, fmt.Errorf("error al comprobar el estado del paso '%s': %w", plannedStep.Step.Name(), err))
			continue // Continuar para marcar los demás como omitidos
		}

		if !decision.Execute {
			plannedStep.Action = planning.ActionSkipCached
			plannedStep.Reason = decision.Reason
			// TODO: Cargar variables desde el estado cacheado si es necesario.
			continue
		}
		plannedStep.Action = planning.ActionExecute

		// 2. Enriquecer variables (versión, commit)
		isTestStepOnly := len(input.Steps) > 0 && i == len(input.Steps)-1 && plannedStep.Step.Name() == "test"
		version, commit, err := o.versionCalculator.CalculateNextVersion(ctx, input.Workspace.Path, isTestStepOnly)
		if err != nil {
			firstError = o.handleStepError(plan, i, fmt.Errorf("error al calcular la versión para el paso '%s': %w", plannedStep.Step.Name(), err))
			continue
		}
		globalVars["var.commit_sha"] = commit.Hash
		globalVars["var.service_version"] = version.Raw

		// 3. Ejecutar el paso
		result, execErr := o.stepExecutor.Execute(ctx, plannedStep.Step, globalVars, input.Workspace.Path)
		if execErr != nil || (result != nil && result.Status == execvos.Failure) {
			errMsg := o.getExecutionErrorMessage(execErr, result)
			firstError = o.handleStepError(plan, i, fmt.Errorf("error al ejecutar el paso '%s': %s", plannedStep.Step.Name(), errMsg))
			continue
		}

		// 4. Actualizar estado y variables globales
		if err := o.stateManager.SaveState(ctx, plannedStep.Step, input.Workspace.Path); err != nil {
			fmt.Printf("ADVERTENCIA: no se pudo guardar el estado para el paso '%s': %v\n", plannedStep.Step.Name(), err)
		}
		for k, v := range result.OutputVars {
			globalVars[k] = v
		}
	}

	return plan, firstError
}

func (o *ExecutionOrchestrator) initializeGlobalVariables(input *OrchestratorInput) execvos.VariableSet {
	return execvos.VariableSet{
		"project.id":   input.Project.ID,
		"project.name": input.Project.Name,
		"project.team": input.Project.Team,
	}
}

func (o *ExecutionOrchestrator) createInitialPlan(steps []*entities.StepDefinition) *planning.ExecutionPlan {
	plannedSteps := make([]*planning.PlannedStep, len(steps))
	for i, step := range steps {
		plannedSteps[i] = &planning.PlannedStep{Step: step, Action: planning.ActionExecute} // Asumimos ejecución por defecto
	}
	return planning.NewExecutionPlan(plannedSteps)
}

func (o *ExecutionOrchestrator) handleStepError(plan *planning.ExecutionPlan, failedIndex int, err error) error {
	plan.Steps[failedIndex].Action = planning.ActionSkipFailedDep
	plan.Steps[failedIndex].Reason = err.Error()

	for j := failedIndex + 1; j < len(plan.Steps); j++ {
		plan.Steps[j].Action = planning.ActionSkipFailedDep
		plan.Steps[j].Reason = fmt.Sprintf("Omitido porque el paso '%s' falló.", plan.Steps[failedIndex].Step.Name())
	}

	// Devuelve el primer error que ocurrió.
	return err
}

func (o *ExecutionOrchestrator) getExecutionErrorMessage(execErr error, result *execvos.ExecutionResult) string {
	if execErr != nil {
		return execErr.Error()
	}
	if result != nil && result.Error != nil {
		return result.Error.Error()
	}
	return "la ejecución del paso falló sin un error explícito"
}
 */