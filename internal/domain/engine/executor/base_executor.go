package executor

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"time"
)

const (
	defaultEmptyTimeout = ""
	defaultRetryDelay   = "1s"
)

type Executor interface {
	Execute(ctx context.Context, step model.Step) error
}

type BaseExecutor struct{}

func NewBaseExecutor() *BaseExecutor {
	return &BaseExecutor{}
}

func (e *BaseExecutor) prepareContext(ctx context.Context, step model.Step) (context.Context, context.CancelFunc) {
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

func (e *BaseExecutor) handleRetry(step model.Step, executionFunc func() error) error {
	if step.Retry == nil {
		return executionFunc()
	}

	var lastError error

	for attemptNum := 0; attemptNum < step.Retry.Attempts; attemptNum++ {
		if attemptNum > 0 {
			delay, err := time.ParseDuration(step.Retry.Delay)
			if err != nil {
				// If delay parsing fails, use a default delay
				delay, _ = time.ParseDuration(defaultRetryDelay)
			}
			time.Sleep(delay)
		}

		err := executionFunc()
		if err == nil {
			// Success, return immediately
			return nil
		}

		lastError = err
	}

	return lastError
}
