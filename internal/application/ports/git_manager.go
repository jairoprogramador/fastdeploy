package ports

import "context"

type GitManager interface {
	IsGit() (bool, error)
	GetCommitHash(ctx context.Context) (string, error)
	ExistChanges(ctx context.Context) (bool, error)
}