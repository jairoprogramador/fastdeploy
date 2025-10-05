package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git/dto"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
)

func EnvironmentsToDomain(dto []dto.EnvironmentDTO) ([]vos.Environment, error) {
	var environments []vos.Environment
	for _, dto := range dto {
		env, err := vos.NewEnvironment(dto.Name, dto.Description, dto.Value)
		if err != nil {
			return []vos.Environment{}, err
		}
		environments = append(environments, env)
	}
	return environments, nil
}