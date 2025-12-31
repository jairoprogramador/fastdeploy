package ports

import "context"

type CommandResultDTO struct {
	Output   string
	ExitCode int
}

type ClonerTemplate interface {
	EnsureCloned(ctx context.Context, repoURL, ref, localPath string) error
	Run(ctx context.Context, command string, workDir string) (*CommandResultDTO, error)
}
