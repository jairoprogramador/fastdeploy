package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git/dto"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
)

func OutputsToDomain(dto []dto.OutputProbeDTO) ([]vos.OutputProbe, error) {
	var outputs []vos.OutputProbe
	for _, dto := range dto {
		output, err := vos.NewOutputProbe(dto.Name, dto.Description, dto.Probe)
		if err != nil {
			return []vos.OutputProbe{}, err
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}