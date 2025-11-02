package ports

import "context"

type GitManager interface {
	IsGit(pathProject string) (bool, error)
	GetCommitHash(ctx context.Context, pathProject string) (string, error)
	ExistChanges(ctx context.Context, pathProject string) (bool, error)
}