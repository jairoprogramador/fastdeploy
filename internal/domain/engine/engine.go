package engine

import (
	"context"
	entity2 "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/executor"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/store"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/validator"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"

	//"deploy/internal/domain/service"
	"fmt"
	"sync"
)

// Error message constants
const (
	errDeploymentValidation  = "deployment file with errors: %v"
	errStepExecution         = "error in step %s: %v"
	errUnsupportedStepType   = "unsupported step type: %s"
	errParallelStepExecution = "error in parallel step %s: %v"
	errMultipleParallelSteps = "errors in parallel steps: %v"
)

// Engine handles the execution of deployment steps
type Engine struct {
	validator     *validator.Validator
	Executors     map[string]executor.Executor
	variableStore *entity2.StoreEntity
	storeService  store.StoreServiceInterface
}

// NewEngine creates a new deployment engine instance
func NewEngine(
	variableStore *entity2.StoreEntity,
	storeService store.StoreServiceInterface,
	validator *validator.Validator,
) *Engine {
	return &Engine{
		validator:     validator,
		Executors:     make(map[string]executor.Executor),
		variableStore: variableStore,
		storeService:  storeService,
	}
}

// Execute runs a deployment process with the given context and configuration
func (e *Engine) Execute(ctx context.Context, deployment *entity2.DeploymentEntity, project *entity.ProjectEntity) error {
	if err := e.validator.Validate(deployment); err != nil {
		return fmt.Errorf(errDeploymentValidation, err)
	}

	globalVars, err := e.storeService.GetVariablesGlobal(ctx, deployment, project)
	if err != nil {
		return err
	}

	e.variableStore.Initialize(globalVars)

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

// executeStep runs a single deployment step or parallel steps
func (e *Engine) executeStep(ctx context.Context, step entity2.Step) error {
	// Handle parallel steps if present
	if len(step.Parallel) > 0 {
		return e.executeParallelSteps(ctx, step.Parallel)
	}

	// Find the appropriate executor for this step type
	executor, exists := e.Executors[step.Type]
	if !exists {
		return fmt.Errorf(errUnsupportedStepType, step.Type)
	}

	// Manage variable scope for this step
	e.variableStore.PushScope(step.Variables)
	defer e.variableStore.PopScope()

	// Execute the step
	return executor.Execute(ctx, step)
}

// executeParallelSteps runs multiple steps concurrently
func (e *Engine) executeParallelSteps(ctx context.Context, steps []entity2.Step) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(steps))

	// Launch each step in its own goroutine
	for _, step := range steps {
		wg.Add(1)
		go func(s entity2.Step) {
			defer wg.Done()

			if err := e.executeStep(ctx, s); err != nil {
				errChan <- fmt.Errorf(errParallelStepExecution, s.Name, err)
			}
		}(step)
	}

	// Wait for all steps to complete
	wg.Wait()
	close(errChan)

	// Collect any errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf(errMultipleParallelSteps, errors)
	}

	return nil
}
