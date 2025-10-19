package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/state/aggregates"
)

type FingerprintRepository interface {
	FindCode() (*aggregates.FingerprintState, error)
	FindStep(stepName string) (*aggregates.FingerprintState, error)
	SaveCode(state *aggregates.FingerprintState) error
	SaveStep(state *aggregates.FingerprintState) error
}
