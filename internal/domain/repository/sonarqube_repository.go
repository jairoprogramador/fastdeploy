package repository

import "deploy/internal/domain/model"
import "context"
import "time"

type SonarqubeRepository interface {
	Add() *model.Response
	WaitSonarqube(ctx context.Context, maxRetries int, interval time.Duration) error
	CreateToken() (string, error)
	ChangePassword() (string, error)
	RevokeToken() error
}
