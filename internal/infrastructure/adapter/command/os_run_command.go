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

type osRunCommand struct {
	fileLogger *logger.FileLogger
}

func NewOsRunCommand(fileLogger *logger.FileLogger) port.RunCommand {
	return &osRunCommand{
		fileLogger: fileLogger,
	}
}

func (r *osRunCommand) Run(ctx context.Context, cmdExec string) result.InfraResult {
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

func (r *osRunCommand) ensureContextTimeout(ctx context.Context) context.Context {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		newCtx, _ := context.WithTimeout(ctx, DefaultCommandTimeout)
		return newCtx
	}
	return ctx
}

func (r *osRunCommand) parseCommand(cmdExec string) (string, []string) {
	parts := strings.Fields(cmdExec)
	command := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}
	return command, args
}

func (r *osRunCommand) prepareCommand(ctx context.Context, command string, args []string, stdout, stderr *bytes.Buffer) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}

func (r *osRunCommand) formatCommand(command string, args []string) string {
	return fmt.Sprintf("Running command: '%s %s'", command, strings.Join(args, " "))
}

func (r *osRunCommand) handleCommandError(ctx context.Context, err error, stdout, stderr *bytes.Buffer) result.InfraResult {
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

func (r *osRunCommand) appendStdMessage(stdout, stderr *bytes.Buffer) {
	stdOutput := strings.TrimSpace(stdout.String())
	if stdOutput != "" {
		r.logError(fmt.Errorf("Standard Output:\n%s", stdOutput))
	}

	stdError := strings.TrimSpace(stderr.String())
	if stdError != "" {
		r.logError(fmt.Errorf("Standard Error:\n%s\n", stdError))
	}
}

func (r *osRunCommand) logError(err error) result.InfraResult {
	if err != nil {
		r.fileLogger.Error(err)
	}
	return result.NewError(err)
}
