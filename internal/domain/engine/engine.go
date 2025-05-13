package engine

import (
	"context"
	"deploy/internal/domain/executor"
	"fmt"
	"sync"

	//"deploy/internal/domain/metrics"
	"deploy/internal/domain/model"
	"deploy/internal/domain/validator"
	"deploy/internal/domain/variable"
	"deploy/internal/domain/service"
	"deploy/internal/infrastructure/repository"
)

type Engine struct {
	validator *validator.DeploymentValidator
	executors map[string]executor.StepExecutor
	//metrics       *metrics.MetricsCollector
	variableStore *variable.VariableStore
	variableService service.VariableServiceInterface
	muEngine      sync.RWMutex
}
func NewEngine() *Engine {
	globalRepo := repository.GetGlobalConfigRepository()
	globalConfigService := service.GetGlobalConfigService(globalRepo)

	projectRepo := repository.GetProjectRepository()
	projectService := service.GetProjectService(projectRepo, globalConfigService)

	variableRepository := repository.GetVariableRepository()
	variableService := service.GetVariableService(projectService, variableRepository)

	e := &Engine{
		validator:     validator.GetDeploymentValidator(),
		executors:     make(map[string]executor.StepExecutor),
		variableStore: variable.GetVariableStore(),
		variableService: variableService,
	}

	//e.RegisterExecutor(validator.TypeCheck, executor.GetCheckExecutor(e.variableStore))
	e.RegisterExecutor(validator.TypeCommand, executor.GetCommandExecutor(e.variableStore))
	e.RegisterExecutor(validator.TypeContainer, executor.GetContainerExecutor(e.variableStore))
	return e
}

func (e *Engine) Execute(ctx context.Context, deployment *model.Deployment) error {
	// Validar el despliegue
	if err := e.validator.Validate(deployment); err != nil {
		return fmt.Errorf("archivo deployment con errores: %v", err)
	}

	// Inicializar variables globales
	e.variableStore.Initialize(deployment.Variables.Global)
	e.variableService.InitializeGlobalDefault(e.variableStore)

	// Ejecutar pasos
	for _, step := range deployment.Steps {
		if err := e.executeStep(ctx, step); err != nil {
			return fmt.Errorf("error en paso %s: %v", step.Name, err)
		}
	}
	return nil
}

func (e *Engine) executeStep(ctx context.Context, step model.Step) error {
	// Verificar dependencias
	/* if err := e.checkDependencies(step); err != nil {
		return err
	} */

	// Ejecutar pasos paralelos si existen
	if len(step.Parallel) > 0 {
		return e.executeParallelSteps(ctx, step.Parallel)
	}

	// Obtener ejecutor para el tipo de paso
	executor, exists := e.GetExecutor(step.Type)
	if !exists {
		return fmt.Errorf("tipo de paso no soportado: %s", step.Type)
	}

	// Preparar variables locales
	e.variableStore.PushScope(step.Variables)
	defer e.variableStore.PopScope()

	// Ejecutar paso
	return executor.Execute(ctx, step)
}

// checkDependencies verifica que todas las dependencias del paso estén satisfechas
/* func (e *Engine) checkDependencies(step model.Step) error {
	for _, dep := range step.DependsOn {
		// Verificar si la dependencia existe en el mapa de pasos ejecutados
		if !e.isStepExecuted(dep) {
			return fmt.Errorf("dependencia no satisfecha: %s", dep)
		}
	}
	return nil
} */

func (e *Engine) executeParallelSteps(ctx context.Context, steps []model.Step) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(steps))

	for _, step := range steps {
		wg.Add(1)
		go func(s model.Step) {
			defer wg.Done()
			if err := e.executeStep(ctx, s); err != nil {
				errChan <- fmt.Errorf("error en paso paralelo %s: %v", s.Name, err)
			}
		}(step)
	}

	// Esperar a que todos los pasos terminen
	wg.Wait()
	close(errChan)

	// Recolectar errores
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores en pasos paralelos: %v", errors)
	}

	return nil
}

// isStepExecuted verifica si un paso ya fue ejecutado
/* func (e *Engine) isStepExecuted(stepName string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.metrics.IsStepExecuted(stepName)
} */

// RegisterExecutor registra un nuevo ejecutor para un tipo de paso
func (e *Engine) RegisterExecutor(stepType string, executor executor.StepExecutor) {
	e.muEngine.Lock()
	defer e.muEngine.Unlock()
	e.executors[stepType] = executor
}

func (e *Engine) GetExecutor(stepType string) (executor.StepExecutor, bool) {
	e.muEngine.RLock() // Bloqueo compartido para lectura
	defer e.muEngine.RUnlock()
	executorRegistered, exists := e.executors[stepType]
	if !exists {
		return executor.GetCommandExecutor(e.variableStore), true
	}
	return executorRegistered, exists
}

// GetMetrics retorna las métricas del despliegue
/* func (e *Engine) GetMetrics() *metrics.DeploymentMetrics {
	return e.metrics.GetMetrics()
} */

/* func (e *Engine) GetVariableStore() *variable.VariableStore {
	return e.variableStore
} */

