package services

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"
)


type FingerprintService interface {
	CalculateCodeFingerprint() (vos.Fingerprint, error)
	CalculateStepFingerprint(stepName string) (vos.Fingerprint, error)
}
