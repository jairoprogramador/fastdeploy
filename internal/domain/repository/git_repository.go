package repository

import "context"

type GitRepository interface {
	GetCommitHash(ctx context.Context) (string, error)
	GetCommitAuthor(ctx context.Context, commitHash string) (string, error)
	GetCommitMessage(ctx context.Context, commitHash string) (string, error)
}