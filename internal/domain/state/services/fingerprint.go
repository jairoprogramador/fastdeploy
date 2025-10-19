package services

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"
)

type FingerprintService interface {
	GenerateFromSource(sourcePath string) (vos.Fingerprint, error)
	GenerateFromStepDefinition(templatePath, stepName string) (vos.Fingerprint, error)
	GenerateFromStepVariables(vars map[string]string) (vos.Fingerprint, error)
}
