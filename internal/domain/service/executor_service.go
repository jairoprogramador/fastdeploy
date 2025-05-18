package service

import (
	"deploy/internal/domain/model"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type ExecutorServiceInterface interface {
	Run(ctx context.Context, command string) (string, error)
}

type DefaultExecutorService struct{
	logStore *model.LogStore
}

var (
	instanceExecutorService *DefaultExecutorService
	onceExecutorService     sync.Once
)

func GetExecutorService() ExecutorServiceInterface {
	onceExecutorService.Do(func() {
		instanceExecutorService = &DefaultExecutorService{
			logStore: model.GetLogStore(),
		}
	})
	return instanceExecutorService
}

func (r *DefaultExecutorService) Run(ctx context.Context, cmdExec string) (string, error) {
	if cmdExec == "" {
		return "", fmt.Errorf("comando vacío")
	}

	var cancel context.CancelFunc
	_, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		ctx, cancel = context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()
	}

	parts := strings.Fields(cmdExec)
	command := parts[0]
	args := parts[1:]

	r.logStore.AddCommand(fmt.Sprintf("%s %s", command, strings.Join(args, " ")))

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("Error ejecutando comando '%s %s': %v",
			command,
			strings.Join(args, " "),
			err,
		)

		switch {
		case ctx.Err() == context.DeadlineExceeded:
			return "", fmt.Errorf("timeout al ejecutar comando: %s", errMsg)
		case ctx.Err() == context.Canceled:
			return "", fmt.Errorf("comando cancelado: %s", errMsg)
		}

		if stdoutBuf.Len() > 0 {
			errMsg += fmt.Sprintf("\nStandard Output:\n%s", stdoutBuf.String())
		}
		if stderrBuf.Len() > 0 {
			errMsg += fmt.Sprintf("\nStandard Error:\n%s", stderrBuf.String())
		}
		return "", fmt.Errorf("%s", errMsg)
	}

	stderrBuf.Reset()
	return strings.TrimSpace(stdoutBuf.String()), nil
}
