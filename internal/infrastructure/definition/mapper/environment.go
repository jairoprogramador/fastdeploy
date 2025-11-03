package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/definition/dto"
)

func EnvironmentsToDomain(dto []dto.EnvironmentDefinitionDTO) ([]vos.EnvironmentDefinition, error) {
	var environments []vos.EnvironmentDefinition
	for _, dto := range dto {
		env, err := vos.NewEnvironmentDefinition(dto.Name, dto.Value)
		if err != nil {
			return []vos.EnvironmentDefinition{}, err
		}
		environments = append(environments, env)
	}
	return environments, nil
}
