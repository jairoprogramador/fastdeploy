package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/definition/dto"
)

func CommandsToDomain(dto []dto.CommandDefinitionDTO) ([]vos.CommandDefinition, error) {
	var commands []vos.CommandDefinition
	for _, d := range dto {
		opts := make([]vos.CommandOption, 0, 4)

		if d.Workdir != "" {
			opts = append(opts, vos.WithWorkdir(d.Workdir))
		}
		if d.TemplateFiles != nil {
			opts = append(opts, vos.WithTemplateFiles(d.TemplateFiles))
		}
		if d.Outputs != nil {
			outputs, err := OutputsToDomain(d.Outputs)
			if err != nil {
				return nil, err
			}
			opts = append(opts, vos.WithOutputs(outputs))
		}

		command, err := vos.NewCommandDefinition(d.Name, d.Cmd, opts...)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, nil
}
