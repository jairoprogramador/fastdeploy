package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git/dto"
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
		if d.Description != "" {
			opts = append(opts, vos.WithDescription(d.Description))
		}

		command, err := vos.NewCommandDefinition(d.Name, d.Cmd, opts...)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, nil
}
