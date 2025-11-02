package services

import (
	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

type FingerprintService interface {
	GenerateFromPath(pathProject string) (vos.Fingerprint, error)
	GenerateFromStepDefinition(pathTemplate string, runParams appDto.RunParams) (vos.Fingerprint, error)
	GenerateFromStepVariables(vars map[string]string) (vos.Fingerprint, error)
}
