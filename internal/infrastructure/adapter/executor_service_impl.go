package adapter

import (
	"bytes"
	"context"
	"deploy/internal/domain/model"
	"deploy/internal/domain/port"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ExecutorServiceImpl struct {
}

func NewExecutorService() port.ExecutorServiceInterface {
	return &ExecutorServiceImpl{}
}

func (r *ExecutorServiceImpl) Run(ctx context.Context, cmdExec string) model.InfrastructureResponse {
	if cmdExec == "" {
		return model.NewErrorResponse(fmt.Errorf("the command is empty"))
	}

	var cancel context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var newCtx context.Context
		newCtx, cancel = context.WithTimeout(ctx, 10*time.Minute)
		ctx = newCtx
		defer cancel()
	}

	parts := strings.Fields(cmdExec)
	command := parts[0]
	args := parts[1:]

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	var message strings.Builder
	message.WriteString(fmt.Sprintf("Running command: '%s %s'\n", command, strings.Join(args, " ")))

	err := cmd.Run()
	if err != nil {
		message.WriteString(fmt.Sprintf("Error running command: %v\n", err))

		if ctxErr := ctx.Err(); ctxErr == context.DeadlineExceeded {
			message.WriteString(fmt.Sprintf("timeout (context deadline exceeded) running command: %v\n", ctx.Err()))
			return model.NewErrorResponseWithDetails(ctx.Err(), message.String())
		} else if ctxErr == context.Canceled {
			message.WriteString(fmt.Sprintf("command canceled (context canceled): %v\n", context.Canceled))
			return model.NewErrorResponseWithDetails(context.Canceled, message.String())
		}

		stdOutput := strings.TrimSpace(stdoutBuf.String())
		stdError := strings.TrimSpace(stderrBuf.String())

		if stdOutput != "" {
			message.WriteString(fmt.Sprintf("Standard Output:\n%s\n", stdOutput))
		}
		if stdError != "" {
			message.WriteString(fmt.Sprintf("Standard Error:\n%s\n", stdError))
		}
		return model.NewErrorResponseWithDetails(err, message.String())
	}

	return model.NewResponseWithDetails(strings.TrimSpace(stdoutBuf.String()), message.String())
}
