package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/deployment/dto"
)

func OutputsToDomain(dto []dto.OutputProbeDTO) ([]vos.Output, error) {
	var outputs []vos.Output
	for _, dto := range dto {
		output, err := vos.NewOutput(dto.Name, dto.Probe)
		if err != nil {
			return []vos.Output{}, err
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}
