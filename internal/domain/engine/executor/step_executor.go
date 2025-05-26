package executor

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/model"
	"time"
)

// Default values as constants
const (
	defaultEmptyTimeout = ""
	defaultRetryDelay   = "1s"
)

// Executor defines the contract for step executors
type Executor interface {
	Execute(ctx context.Context, step model.Step) error
}

// StepExecutor provides common functionality for all step executors
type StepExecutor struct{}

// NewStepExecutor creates a new instance of StepExecutor
func NewStepExecutor() *StepExecutor {
	return &StepExecutor{}
}

// prepareContext creates a context with timeout if specified in the step
// Returns the original context and a no-op cancel function if:
// - No timeout is specified
// - The timeout string cannot be parsed
func (e *StepExecutor) prepareContext(ctx context.Context, step model.Step) (context.Context, context.CancelFunc) {
	// Return early if no timeout specified
	if step.Timeout == defaultEmptyTimeout {
		return ctx, func() {}
	}

	timeout, err := time.ParseDuration(step.Timeout)
	if err != nil {
		// Log or handle the error appropriately in a real implementation
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, timeout)
}

// handleRetry executes the provided function with retry logic based on step configuration
// If retry is not configured, executes the function once
// If retry is configured, attempts execution up to the specified number of times
func (e *StepExecutor) handleRetry(step model.Step, executionFunc func() error) error {
	// Execute once if no retry configuration
	if step.Retry == nil {
		return executionFunc()
	}

	var lastError error

	// Try execution up to the configured number of attempts
	for attemptNum := 0; attemptNum < step.Retry.Attempts; attemptNum++ {
		// Apply delay between retries, but not before the first attempt
		if attemptNum > 0 {
			delay, err := time.ParseDuration(step.Retry.Delay)
			if err != nil {
				// If delay parsing fails, use a default delay
				delay, _ = time.ParseDuration(defaultRetryDelay)
			}
			time.Sleep(delay)
		}

		// Execute the function
		err := executionFunc()
		if err == nil {
			// Success, return immediately
			return nil
		}

		// Store the error for potential return
		lastError = err
	}

	// All attempts failed, return the last error
	return lastError
}
