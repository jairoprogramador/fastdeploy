package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/definition/dto"
)

func OutputsToDomain(dto []dto.OutputDefinitionDTO) ([]vos.OutputDefinition, error) {
	var outputs []vos.OutputDefinition
	for _, dto := range dto {
		output, err := vos.NewOutputDefinition(dto.Name, dto.Probe)
		if err != nil {
			return []vos.OutputDefinition{}, err
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}
