package ports

import appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"

type VariablesRepository interface {
	FindByStep(namesRequest appDto.NamesParams, runParams appDto.RunParams) (map[string]string, error)
	SaveByStep(namesRequest appDto.NamesParams, runParams appDto.RunParams, vars map[string]string) error
}
