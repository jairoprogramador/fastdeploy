package executor

import (
/* 	"context"
	"deploy/internal/domain/model"
	"fmt" */
	"net/http"
)

type HTTPExecutor struct {
	BaseExecutor
	client *http.Client
}

/* func (e *HTTPExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.prepareContext(ctx, step)
	defer cancel()

	return e.handleRetry(ctx, step, func() error {
		req, err := http.NewRequestWithContext(ctx, step.Method, step.URL, nil)
		if err != nil {
			return err
		}

		resp, err := e.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != step.ExpectedStatus {
			return fmt.Errorf("código de estado inesperado: %d", resp.StatusCode)
		}

		return nil
	})
} */
