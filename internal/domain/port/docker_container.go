package port

import (
	"context"
	"deploy/internal/domain/model"
)

type DockerContainer interface {
	UpBuild(ctx context.Context) model.InfrastructureResponse
	Up(ctx context.Context) model.InfrastructureResponse
	Down(ctx context.Context) model.InfrastructureResponse
	Exists(ctx context.Context, commitHash, version string) model.InfrastructureResponse
	GetURLs(ctx context.Context, commitHash, version string) model.InfrastructureResponse
	Create(ctx context.Context) error
	Start(ctx context.Context) error
}
