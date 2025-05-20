package service

import (
	"context"
)

type ExecutorServiceInterface interface {
	Run(ctx context.Context, command string) (string, error)
}
