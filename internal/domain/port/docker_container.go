package port

import (
	"context"
	"deploy/internal/domain/model"
)

type DockerContainer interface {
	Up(ctx context.Context) model.InfrastructureResponse
	Exists(ctx context.Context, commitHash, version string) model.InfrastructureResponse
	GetURLsUp(ctx context.Context, commitHash, version string) model.InfrastructureResponse
	Start(ctx context.Context) error
}
