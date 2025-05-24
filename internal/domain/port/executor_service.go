package port

import (
	"context"
	"deploy/internal/domain/model"
)

type ExecutorServiceInterface interface {
	Run(ctx context.Context, command string) model.InfrastructureResponse
}
