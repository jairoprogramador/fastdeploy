package command

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"os/exec"
	"strings"
	"time"
)

const (
	DefaultCommandTimeout = 10 * time.Minute

	errCommandEmpty    = "the command is empty"
	errTimeoutExceeded = "timeout (context deadline exceeded): %v"
	errCommandCanceled = "command canceled: %v"
)

type commandAdapter struct {
	fileLogger *logger.FileLogger
}

func NewCommandAdapter(fileLogger *logger.FileLogger) port.CommandPort {
	return &commandAdapter{
		fileLogger: fileLogger,
	}
}

func (r *commandAdapter) Run(ctx context.Context, cmdExec string) result.InfraResult {
	if cmdExec == "" {
		return r.logError(fmt.Errorf(errCommandEmpty))
	}

	ctx = r.ensureContextTimeout(ctx)
	command, args := r.parseCommand(cmdExec)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := r.prepareCommand(ctx, command, args, &stdoutBuf, &stderrBuf)

	formatCmd := r.formatCommand(command, args)
	r.fileLogger.Info(formatCmd)

	err := cmd.Run()
	if err != nil {
		return r.handleCommandError(ctx, err, &stdoutBuf, &stderrBuf)
	}

	return result.NewResult(strings.TrimSpace(stdoutBuf.String()))
}

func (r *commandAdapter) ensureContextTimeout(ctx context.Context) context.Context {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		newCtx, _ := context.WithTimeout(ctx, DefaultCommandTimeout)
		return newCtx
	}
	return ctx
}

func (r *commandAdapter) parseCommand(cmdExec string) (string, []string) {
	parts := strings.Fields(cmdExec)
	command := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}
	return command, args
}

func (r *commandAdapter) prepareCommand(ctx context.Context, command string, args []string, stdout, stderr *bytes.Buffer) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}

func (r *commandAdapter) formatCommand(command string, args []string) string {
	return fmt.Sprintf("Running command: '%s %s'", command, strings.Join(args, " "))
}

func (r *commandAdapter) handleCommandError(ctx context.Context, err error, stdout, stderr *bytes.Buffer) result.InfraResult {
	r.appendStdMessage(stdout, stderr)

	if ctxErr := ctx.Err(); ctxErr != nil {
		if errors.Is(ctxErr, context.DeadlineExceeded) {
			return r.logError(fmt.Errorf(errTimeoutExceeded, ctxErr))
		}

		if errors.Is(ctxErr, context.Canceled) {
			return r.logError(fmt.Errorf(errCommandCanceled, ctxErr))
		}
	}
	return r.logError(err)
}

func (r *commandAdapter) appendStdMessage(stdout, stderr *bytes.Buffer) {
	stdOutput := strings.TrimSpace(stdout.String())
	if stdOutput != "" {
		r.logError(fmt.Errorf("Standard Output:\n%s", stdOutput))
	}

	stdError := strings.TrimSpace(stderr.String())
	if stdError != "" {
		r.logError(fmt.Errorf("Standard Error:\n%s\n", stdError))
	}
}

func (r *commandAdapter) logError(err error) result.InfraResult {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return result.NewError(err)
}
