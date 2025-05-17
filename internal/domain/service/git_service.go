package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type GitServiceInterface interface {
	GetCommitHash(ctx context.Context) (string, error)
	GetCommitAuthor(ctx context.Context, commitHash string) (string, error)
	GetCommitMessage(ctx context.Context, commitHash string) (string, error)
}

type GitService struct {
	executorService ExecutorServiceInterface
}

var (
	instanceGitService     *GitService
	instanceOnceGitService sync.Once
)

func GetGitService() GitServiceInterface {
	instanceOnceGitService.Do(func() {
		instanceGitService = &GitService{
			executorService: GetExecutorService(),
		}
	})
	return instanceGitService
}

func (s *GitService) GetCommitHash(ctx context.Context) (string, error) {
	command := "git rev-parse HEAD"
	return s.getResult(ctx, command)
}

func (s *GitService) GetCommitAuthor(ctx context.Context, commitHash string) (string, error) {
	command := fmt.Sprintf("git show -s --format=%%an<%%ae> %s", commitHash)
	return s.getResult(ctx, command)
}

func (s *GitService) GetCommitMessage(ctx context.Context, commitHash string) (string, error) {
	command := fmt.Sprintf("git show -s --format=%%s %s", commitHash)
	return s.getResult(ctx, command)
}

func (s *GitService) getResult(ctx context.Context, command string) (string, error) {
	result, err := s.executorService.Run(ctx, command)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}
