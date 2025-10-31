package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/template/dto"
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
