package port

import (
	"context"
	"deploy/internal/domain/model"
)

type RunCommand interface {
	Run(ctx context.Context, command string) model.InfrastructureResponse
}
