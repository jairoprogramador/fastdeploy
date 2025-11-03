package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/definition/dto"
)

func VariablesToDomain(dtos []dto.VariableDefinitionDTO) ([]vos.VariableDefinition, error) {
	variables := make([]vos.VariableDefinition, 0, len(dtos))
	for _, dto := range dtos {
		variable, err := vos.NewVariableDefinition(dto.Name, dto.Value)
		if err != nil {
			return nil, err
		}
		variables = append(variables, variable)
	}
	return variables, nil
}
