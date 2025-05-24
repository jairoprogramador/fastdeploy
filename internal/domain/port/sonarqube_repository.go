package port

import (
	"context"
	"deploy/internal/domain/model"
	"time"
)

type SonarqubeRepository interface {
	Add() model.InfrastructureResponse
	WaitSonarqube(ctx context.Context, maxRetries int, interval time.Duration) model.InfrastructureResponse
	CreateToken(projectKey string) model.InfrastructureResponse
	ChangePassword() model.InfrastructureResponse
	GetQualityGateStatus(projectKey string) model.InfrastructureResponse
	RevokeToken() model.InfrastructureResponse
}
