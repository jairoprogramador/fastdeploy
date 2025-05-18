package executor

import (
	"context"
	"deploy/internal/domain/model"
	"time"
	"sync"
)

type StepExecutor interface {
	Execute(ctx context.Context, step model.Step) (string, error)
}

type BaseExecutor struct {}

var (
	instanceBaseExecutor *BaseExecutor
	onceBaseExecutor     sync.Once
)

func GetBaseExecutor() *BaseExecutor {
	onceBaseExecutor.Do(func() {
		instanceBaseExecutor = &BaseExecutor{}
	})
	return instanceBaseExecutor
}

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

func (e *BaseExecutor) handleRetry(step model.Step, fn func() (string, error)) (string, error) {
	if step.Retry == nil {
		return fn()
	}

	var lastErr error
	for i := 0; i < step.Retry.Attempts; i++ {
		if i > 0 {
			delay, _ := time.ParseDuration(step.Retry.Delay)
			time.Sleep(delay)
		}

		if message, err := fn(); err == nil {
			return message, nil
		} else {
			lastErr = err
		}
	}
	return "", lastErr
}
