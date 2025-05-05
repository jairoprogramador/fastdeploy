package repository

import (
	"context"
	"deploy/internal/domain/model"
	"time"
)

type SonarqubeRepository interface {
	Add() *model.Response
	WaitSonarqube(ctx context.Context, maxRetries int, interval time.Duration) error
	CreateToken(projectKey string) (string, error)
	ChangePassword() (string, error)
	GetQualityGateStatus(projectKey string) (string, error)
	RevokeToken() error
}
