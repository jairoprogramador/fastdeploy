package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/deployment/dto"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"
)

func EnvironmentsToDomain(dto []dto.EnvironmentDTO) ([]vos.Environment, error) {
	var environments []vos.Environment
	for _, dto := range dto {
		env, err := vos.NewEnvironment(dto.Name, dto.Value)
		if err != nil {
			return []vos.Environment{}, err
		}
		environments = append(environments, env)
	}
	return environments, nil
}