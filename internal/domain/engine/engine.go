package engine

import (
	"context"
	"fmt"
	"sync"
	"deploy/internal/domain/executor"
	"deploy/internal/domain/model"
	"deploy/internal/domain/validator"
	"deploy/internal/domain/variable"
	"deploy/internal/domain/service"
	"deploy/internal/domain/router"
)

type Engine struct {
	validator *validator.DeploymentValidator
	executors map[string]executor.StepExecutor
	variableStore *variable.VariableStore
	storeService service.StoreServiceInterface
	muEngine      sync.RWMutex
}
func NewEngine(
	variableStore *variable.VariableStore,
	storeService service.StoreServiceInterface) *Engine {

	e := &Engine{
		validator:     validator.GetDeploymentValidator(),
		executors:     make(map[string]executor.StepExecutor),
		variableStore: variableStore,
		storeService:  storeService,
	}
	return e
}

func (e *Engine) Execute(ctx context.Context, deployment *model.Deployment) error {
	if err := e.validator.Validate(deployment); err != nil {
		return fmt.Errorf("archivo deployment con errores: %v", err)
	}

	variablesGlobal, err := e.storeService.GetVariablesGlobal(ctx, deployment)
	if err != nil {
		return fmt.Errorf("error al obtener variables globales: %v", err)
	}

	e.variableStore.Initialize(variablesGlobal)

	if deployment.HasType(validator.TypeContainer) {
		if e.tryContainerUp(ctx, e.variableStore) {
			return nil
		}
	}

	for _, step := range deployment.Steps {
		if err := e.executeStep(ctx, step); err != nil {
			return fmt.Errorf("error en paso %s: %v", step.Name, err)
		}
	}

	return nil
}

func (e *Engine) tryContainerUp(ctx context.Context, variableStore *variable.VariableStore) bool {
	dockerService := service.GetDockerService()
	exists, _ := dockerService.ExistsContainer(ctx, variableStore)
	if exists {
		pathDockerCompose := router.GetRouter().GetFullPathDockerCompose()
		err := dockerService.DockerComposeUp(ctx, pathDockerCompose, variableStore)
		return err == nil
	}
	return false
}

func (e *Engine) executeStep(ctx context.Context, step model.Step) error {
	if len(step.Parallel) > 0 {
		return e.executeParallelSteps(ctx, step.Parallel)
	}

	executor, exists := e.getExecutor(step.Type)
	if !exists {
		return fmt.Errorf("tipo de paso no soportado: %s", step.Type)
	}

	e.variableStore.PushScope(step.Variables)
	defer e.variableStore.PopScope()

	return executor.Execute(ctx, step)
}

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

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores en pasos paralelos: %v", errors)
	}

	return nil
}

func (e *Engine) RegisterExecutor(stepType string, executor executor.StepExecutor) {
	e.muEngine.Lock()
	defer e.muEngine.Unlock()
	e.executors[stepType] = executor
}

func (e *Engine) getExecutor(stepType string) (executor.StepExecutor, bool) {
	e.muEngine.RLock()
	defer e.muEngine.RUnlock()
	executorRegistered, exists := e.executors[stepType]
	if !exists {
		return executor.GetCommandExecutor(e.variableStore), true
	}
	return executorRegistered, exists
}
