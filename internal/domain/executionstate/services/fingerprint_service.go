package services

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"
)


type FingerprintService interface {
	CalculateCodeFingerprint(ctx context.Context, pathProject string) (vos.Fingerprint, error)
	CalculateEnvironmentFingerprint(ctx context.Context, stepName string, pathRepository string) (vos.Fingerprint, error)
}
