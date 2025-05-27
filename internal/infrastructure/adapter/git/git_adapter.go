package git

import (
	"context"
	"errors"
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"strings"
)

const (
	cmdGetHash    = "git rev-parse HEAD"
	cmdShowFormat = "git show -s --format=%s %s"

	errCommitHashEmpty = "commit hash cannot be empty"
	errCommitHashShort = "commit hash %s is too short (minimum 7 characters)"
)

type gitAdapter struct {
	commandPort port.CommandPort
	fileLogger  *logger.FileLogger
}

func NewGitAdapter(
	commandPort port.CommandPort,
	fileLogger *logger.FileLogger,
) port.GitPort {
	return &gitAdapter{
		commandPort: commandPort,
		fileLogger:  fileLogger,
	}
}

func (git *gitAdapter) GetHash(ctx context.Context) result.InfraResult {
	return git.executeCommand(ctx, cmdGetHash)
}

func (git *gitAdapter) GetAuthor(ctx context.Context, commitHash string) result.InfraResult {
	if err := validateCommitHash(commitHash); err != nil {
		return git.logError(err)
	}

	authorFormatCmd := fmt.Sprintf(cmdShowFormat, "%%an<%%ae>", commitHash)
	return git.executeCommand(ctx, authorFormatCmd)
}

func (git *gitAdapter) GetMessage(ctx context.Context, commitHash string) result.InfraResult {
	if err := validateCommitHash(commitHash); err != nil {
		return git.logError(err)
	}

	messageFormatCmd := fmt.Sprintf(cmdShowFormat, "%%s", commitHash)
	return git.executeCommand(ctx, messageFormatCmd)
}

func (git *gitAdapter) executeCommand(ctx context.Context, command string) result.InfraResult {
	return git.commandPort.Run(ctx, command)
}

func validateCommitHash(hash string) error {
	if hash == "" {
		return errors.New(errCommitHashEmpty)
	}

	hash = strings.TrimSpace(hash)
	if len(hash) < 7 {
		return fmt.Errorf(errCommitHashShort, hash)
	}

	return nil
}

func (git *gitAdapter) logError(err error) result.InfraResult {
	if err != nil {
		git.fileLogger.Error(err)
	}
	return result.NewError(err)
}
