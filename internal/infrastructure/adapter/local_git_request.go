package adapter

import (
	"context"
	"deploy/internal/domain/model"
	"deploy/internal/domain/port"
	"errors"
	"fmt"
	"strings"
)

// Git command constants to avoid magic strings
const (
	cmdGetHash    = "git rev-parse HEAD"
	cmdShowFormat = "git show -s --format=%s %s"
)

// localGitRequest implements the GitRequest interface using local git commands
type localGitRequest struct {
	commandRunner port.RunCommand
}

// NewLocalGitRequest creates a new instance of GitRequest that uses local git commands
func NewLocalGitRequest(commandRunner port.RunCommand) port.GitRequest {
	return &localGitRequest{
		commandRunner: commandRunner,
	}
}

// GetHash retrieves the current commit hash from the git repository
func (git *localGitRequest) GetHash(ctx context.Context) model.InfrastructureResponse {
	return git.executeCommand(ctx, cmdGetHash)
}

// GetAuthor retrieves the author name and email for the specified commit
func (git *localGitRequest) GetAuthor(ctx context.Context, commitHash string) model.InfrastructureResponse {
	if err := validateCommitHash(commitHash); err != nil {
		return model.NewErrorResponse(err)
	}

	authorFormatCmd := fmt.Sprintf(cmdShowFormat, "%%an<%%ae>", commitHash)
	return git.executeCommand(ctx, authorFormatCmd)
}

// GetMessage retrieves the commit message for the specified commit
func (git *localGitRequest) GetMessage(ctx context.Context, commitHash string) model.InfrastructureResponse {
	if err := validateCommitHash(commitHash); err != nil {
		return model.NewErrorResponse(err)
	}

	messageFormatCmd := fmt.Sprintf(cmdShowFormat, "%%s", commitHash)
	return git.executeCommand(ctx, messageFormatCmd)
}

// executeCommand runs a git command and returns its result
func (git *localGitRequest) executeCommand(ctx context.Context, command string) model.InfrastructureResponse {
	return git.commandRunner.Run(ctx, command)
}

// validateCommitHash ensures the commit hash is not empty and has a valid format
func validateCommitHash(hash string) error {
	if hash == "" {
		return errors.New("commit hash cannot be empty")
	}

	hash = strings.TrimSpace(hash)
	if len(hash) < 7 {
		return fmt.Errorf("commit hash '%s' is too short (minimum 7 characters)", hash)
	}

	return nil
}
