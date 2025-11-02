package ports

import appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"

type VariablesRepository interface {
	FindByStepName(namesRequest appDto.NamesParams, runParams appDto.RunParams) (map[string]string, error)
	Save(namesRequest appDto.NamesParams, runParams appDto.RunParams, vars map[string]string) error
}
