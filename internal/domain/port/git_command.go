package port

import (
	"context"
	"deploy/internal/domain/model"
)

type GitCommand interface {
	GetHash(ctx context.Context) model.InfrastructureResponse
	GetAuthor(ctx context.Context, commitHash string) model.InfrastructureResponse
	GetMessage(ctx context.Context, commitHash string) model.InfrastructureResponse
}
