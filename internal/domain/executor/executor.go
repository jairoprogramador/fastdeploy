package executor

import (
	"context"
	"deploy/internal/domain/model"
	"time"
)

type StepExecutor interface {
	Execute(ctx context.Context, step model.Step) error
}

type BaseExecutor struct {}

func (e *BaseExecutor) prepareContext(ctx context.Context, step model.Step) (context.Context, context.CancelFunc) {
	if step.Timeout == "" {
		return ctx, func() {}
	}

	timeout, err := time.ParseDuration(step.Timeout)
	if err != nil {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, timeout)
}

func (e *BaseExecutor) handleRetry(step model.Step, fn func() error) error {
	if step.Retry == nil {
		return fn()
	}

	var lastErr error
	for i := 0; i < step.Retry.Attempts; i++ {
		if i > 0 {
			delay, _ := time.ParseDuration(step.Retry.Delay)
			time.Sleep(delay)
		}

		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	return lastErr
}
