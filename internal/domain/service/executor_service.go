package service

import (
	"context"
	"deploy/internal/infrastructure/tools"
	"fmt"
	"strings"
	"sync"
)

type ExecutorServiceInterface interface {
	Run(ctx context.Context, command string) (string, error)
}

type DefaultExecutorService struct{}

var (
	instanceExecutorService *DefaultExecutorService
	onceExecutorService     sync.Once
)

func GetExecutorService() ExecutorServiceInterface {
	onceExecutorService.Do(func() {
		instanceExecutorService = &DefaultExecutorService{}
	})
	return instanceExecutorService
}

func (r *DefaultExecutorService) Run(ctx context.Context, command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("comando vacío")
	}

	cmd := parts[0]
	args := parts[1:]

	return tools.ExecuteCommandWithContext(ctx, cmd, args...)
}
