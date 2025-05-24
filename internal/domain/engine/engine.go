package engine

import (
	"context"
	"deploy/internal/domain/engine/executor"
	"deploy/internal/domain/engine/validator"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/service"
	"fmt"
	"sync"
)

type Engine struct {
	validator     *validator.DeploymentValidator
	Executors     map[string]executor.StepExecutorInterface
	variableStore *model.VariableStore
	storeService  service.StoreServiceInterface
	logger        *logger.Logger
}

func NewEngine(
	variableStore *model.VariableStore,
	storeService service.StoreServiceInterface,
	logger *logger.Logger,
	validator *validator.DeploymentValidator,
) *Engine {
	e := &Engine{
		validator:     validator,
		Executors:     make(map[string]executor.StepExecutorInterface),
		variableStore: variableStore,
		storeService:  storeService,
		logger:        logger,
	}
	return e
}

func (e *Engine) Execute(ctx context.Context, deployment *model.DeploymentEntity, project *model.ProjectEntity) error {
	if err := e.validator.Validate(deployment); err != nil {
		return fmt.Errorf("archivo deployment con errores: %v", err)
	}

	variablesGlobal, err := e.storeService.GetVariablesGlobal(ctx, deployment, project)
	if err != nil {
		return err
	}

	e.variableStore.Initialize(variablesGlobal)

	for _, step := range deployment.Steps {
		e.logger.Info(step.Name)
		err := e.executeStep(ctx, step)
		if err != nil {
			return fmt.Errorf("error en paso %s: %v", step.Name, err)
		} else {
			if step.Then == validator.ThenFinish {
				break
			}
		}
	}

	return nil
}

func (e *Engine) executeStep(ctx context.Context, step model.Step) error {
	if len(step.Parallel) > 0 {
		return e.executeParallelSteps(ctx, step.Parallel)
	}

	executor, exists := e.Executors[step.Type]
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
