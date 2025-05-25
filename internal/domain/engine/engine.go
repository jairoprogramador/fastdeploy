package engine

import (
	"context"
	"deploy/internal/domain/engine/executor"
	engineModel "deploy/internal/domain/engine/model"
	"deploy/internal/domain/engine/validator"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/service"
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
	variableStore *engineModel.VariableStore
	storeService  service.StoreServiceInterface
	logger        *logger.Logger
}

// NewEngine creates a new deployment engine instance
func NewEngine(
	variableStore *engineModel.VariableStore,
	storeService service.StoreServiceInterface,
	logger *logger.Logger,
	validator *validator.Validator,
) *Engine {
	return &Engine{
		validator:     validator,
		Executors:     make(map[string]executor.Executor),
		variableStore: variableStore,
		storeService:  storeService,
		logger:        logger,
	}
}

// Execute runs a deployment process with the given context and configuration
func (e *Engine) Execute(ctx context.Context, deployment *engineModel.DeploymentEntity, project *model.ProjectEntity) error {
	// Validate deployment configuration
	if err := e.validator.Validate(deployment); err != nil {
		return fmt.Errorf(errDeploymentValidation, err)
	}

	// Get global variables
	globalVars, err := e.storeService.GetVariablesGlobal(ctx, deployment, project)
	if err != nil {
		return err
	}

	// Initialize variable store
	e.variableStore.Initialize(globalVars)

	// Execute each step sequentially
	for _, step := range deployment.Steps {
		e.logger.Info(step.Name)

		if err := e.executeStep(ctx, step); err != nil {
			return fmt.Errorf(errStepExecution, step.Name, err)
		}

		// Check if we should finish after this step
		if step.Then == validator.ThenFinish {
			break
		}
	}

	return nil
}

// executeStep runs a single deployment step or parallel steps
func (e *Engine) executeStep(ctx context.Context, step engineModel.Step) error {
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
func (e *Engine) executeParallelSteps(ctx context.Context, steps []engineModel.Step) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(steps))

	// Launch each step in its own goroutine
	for _, step := range steps {
		wg.Add(1)
		go func(s engineModel.Step) {
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

	// Return combined error if any steps failed
	if len(errors) > 0 {
		return fmt.Errorf(errMultipleParallelSteps, errors)
	}

	return nil
}
