package ports

import (
	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
)

type FingerprintRepository interface {
	FindCode(namesRequest appDto.NamesParams) (*aggregates.FingerprintState, error)
	FindStep(namesRequest appDto.NamesParams, orderRequest appDto.RunParams) (*aggregates.FingerprintState, error)
	SaveCode(namesRequest appDto.NamesParams, state *aggregates.FingerprintState) error
	SaveStep(namesRequest appDto.NamesParams, orderRequest appDto.RunParams, state *aggregates.FingerprintState) error
}
