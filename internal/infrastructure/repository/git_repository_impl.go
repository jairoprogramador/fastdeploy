package repository

import (
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/tools"
	"sync"
	"strings"
	"context"
)

type GitRepositoryImpl struct{}

var (
	instanceGitRepository     repository.GitRepository
	instanceOnceGitRepository sync.Once
)

func GetGitRepository() repository.GitRepository {
	instanceOnceGitRepository.Do(func() {
		instanceGitRepository = &GitRepositoryImpl{}
	})
	return instanceGitRepository
}

func (s *GitRepositoryImpl) GetCommitHash(ctx context.Context) (string, error) {
	commitHash, err := tools.ExecuteCommandWithContext(ctx, "git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(commitHash), nil
}

func (s *GitRepositoryImpl) GetCommitAuthor(ctx context.Context, commitHash string) (string, error) {
	commitAuthor, err := tools.ExecuteCommandWithContext(ctx, "git", "show", "-s", "--format=%an <%ae>", commitHash)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(commitAuthor), nil
}

func (s *GitRepositoryImpl) GetCommitMessage(ctx context.Context, commitHash string) (string, error) {
	commitMessage, err := tools.ExecuteCommandWithContext(ctx, "git", "show", "-s", "--format=%s", commitHash)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(commitMessage), nil
}
