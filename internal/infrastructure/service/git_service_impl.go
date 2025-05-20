package service

import (
	"context"
	"deploy/internal/domain/service"
	"fmt"
	"strings"
)

type GitServiceImpl struct {
	executorService service.ExecutorServiceInterface
}

func NewGitServiceImpl(executorService service.ExecutorServiceInterface) service.GitServiceInterface {
	return &GitServiceImpl{
		executorService: executorService,
	}
}

func (s *GitServiceImpl) GetCommitHash(ctx context.Context) (string, error) {
	command := "git rev-parse HEAD"
	return s.getResult(ctx, command)
}

func (s *GitServiceImpl) GetCommitAuthor(ctx context.Context, commitHash string) (string, error) {
	command := fmt.Sprintf("git show -s --format=%%an<%%ae> %s", commitHash)
	return s.getResult(ctx, command)
}

func (s *GitServiceImpl) GetCommitMessage(ctx context.Context, commitHash string) (string, error) {
	command := fmt.Sprintf("git show -s --format=%%s %s", commitHash)
	return s.getResult(ctx, command)
}

func (s *GitServiceImpl) getResult(ctx context.Context, command string) (string, error) {
	result, err := s.executorService.Run(ctx, command)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}
