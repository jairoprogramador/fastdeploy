package service

import (
	"bytes"
	"context"
	"deploy/internal/domain/model"
	"deploy/internal/domain/service"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ExecutorServiceImpl struct {
	logStore *model.LogStore
}

func NewExecutorService(logStore *model.LogStore) service.ExecutorServiceInterface {
	return &ExecutorServiceImpl{
		logStore: logStore,
	}
}

func (r *ExecutorServiceImpl) Run(ctx context.Context, cmdExec string) (string, error) {
	if cmdExec == "" {
		return "", fmt.Errorf("comando vacío")
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

	if r.logStore != nil {
		r.logStore.AddCommand(fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		var errMsgBuilder strings.Builder
		errMsgBuilder.WriteString(fmt.Sprintf("Error ejecutando comando '%s %s': %v", command, strings.Join(args, " "), err))

		// Chequear si el error del contexto es la causa raíz, no solo el error de cmd.Run()
		if ctxErr := ctx.Err(); ctxErr == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout (context deadline exceeded) al ejecutar comando: %s", errMsgBuilder.String())
		} else if ctxErr == context.Canceled {
			return "", fmt.Errorf("comando cancelado (context canceled): %s", errMsgBuilder.String())
		}

		stdOutput := strings.TrimSpace(stdoutBuf.String())
		stdError := strings.TrimSpace(stderrBuf.String())

		if stdOutput != "" {
			errMsgBuilder.WriteString(fmt.Sprintf("\nStandard Output:\n%s", stdOutput))
		}
		if stdError != "" {
			errMsgBuilder.WriteString(fmt.Sprintf("\nStandard Error:\n%s", stdError))
		}
		return "", fmt.Errorf("%s", errMsgBuilder.String())
	}

	return strings.TrimSpace(stdoutBuf.String()), nil
}
