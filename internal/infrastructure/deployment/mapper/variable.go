package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/dto"
)

func VariablesToDomain(dtos []dto.VariableDTO) ([]vos.Variable, error) {
	variables := make([]vos.Variable, 0, len(dtos))
	for _, dto := range dtos {
		variable, err := vos.NewVariable(dto.Name, dto.Value)
		if err != nil {
			return nil, err
		}
		variables = append(variables, variable)
	}
	return variables, nil
}
