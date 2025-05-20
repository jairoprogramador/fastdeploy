package service

import (
	"context"
)

type GitServiceInterface interface {
	GetCommitHash(ctx context.Context) (string, error)
	GetCommitAuthor(ctx context.Context, commitHash string) (string, error)
	GetCommitMessage(ctx context.Context, commitHash string) (string, error)
}
