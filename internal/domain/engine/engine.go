package engine

import (
	"context"
	"fmt"
	modelDeploy "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/executor"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
	"sync"
)

const (
	errDeploymentValidation  = "deployment file with errors: %v"
	errStepExecution         = "error in step %s: %v"
	errUnsupportedStepType   = "unsupported step type: %s"
	errParallelStepExecution = "error in parallel step %s: %v"
	errMultipleParallelSteps = "errors in parallel steps: %v"
)

type Engine struct {
	validator *validator.Validator
	executors map[string]executor.Executor
	variables *modelDeploy.StoreEntity
}

func NewEngine(
	storeVariable *modelDeploy.StoreEntity,
	validator *validator.Validator,
) *Engine {
	return &Engine{
		validator: validator,
		executors: make(map[string]executor.Executor),
		variables: storeVariable,
	}
}

func (e *Engine) Execute(ctx context.Context, deployment *modelDeploy.DeploymentEntity) error {
	if err := e.validator.Validate(deployment); err != nil {
		return fmt.Errorf(errDeploymentValidation, err)
	}

	for _, step := range deployment.Steps {

		if err := e.executeStep(ctx, step); err != nil {
			return fmt.Errorf(errStepExecution, step.Name, err)
		}

		if step.Then == validator.ThenFinish {
			break
		}
	}

	return nil
}

func (e *Engine) executeStep(ctx context.Context, step modelDeploy.Step) error {
	if len(step.Parallel) > 0 {
		return e.executeParallelSteps(ctx, step.Parallel)
	}

	executor, exists := e.executors[step.Type]
	if !exists {
		return fmt.Errorf(errUnsupportedStepType, step.Type)
	}

	e.variables.PushScope(step.Variables)
	defer e.variables.PopScope()

	return executor.Execute(ctx, step)
}

func (e *Engine) executeParallelSteps(ctx context.Context, steps []modelDeploy.Step) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(steps))

	for _, step := range steps {
		wg.Add(1)
		go func(s modelDeploy.Step) {
			defer wg.Done()

			if err := e.executeStep(ctx, s); err != nil {
				errChan <- fmt.Errorf(errParallelStepExecution, s.Name, err)
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
		return fmt.Errorf(errMultipleParallelSteps, errors)
	}

	return nil
}

func (e *Engine) AddExecutor(stepType modelDeploy.TypeStep, executor executor.Executor) {
	e.executors[string(stepType)] = executor
}
