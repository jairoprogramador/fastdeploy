package adapter

import (
	"context"
	"deploy/internal/domain/model"
	"deploy/internal/domain/port"
	"fmt"
)

type LocalGitCommand struct {
	executorService port.ExecutorServiceInterface
}

func NewLocalGitCommand(executorService port.ExecutorServiceInterface) port.GitCommand {
	return &LocalGitCommand{
		executorService: executorService,
	}
}

func (s *LocalGitCommand) GetHash(ctx context.Context) model.InfrastructureResponse {
	command := "git rev-parse HEAD"
	return s.getResult(ctx, command)
}

func (s *LocalGitCommand) GetAuthor(ctx context.Context, commitHash string) model.InfrastructureResponse {
	command := fmt.Sprintf("git show -s --format=%%an<%%ae> %s", commitHash)
	return s.getResult(ctx, command)
}

func (s *LocalGitCommand) GetMessage(ctx context.Context, commitHash string) model.InfrastructureResponse {
	command := fmt.Sprintf("git show -s --format=%%s %s", commitHash)
	return s.getResult(ctx, command)
}

func (s *LocalGitCommand) getResult(ctx context.Context, command string) model.InfrastructureResponse {
	return s.executorService.Run(ctx, command)
}
