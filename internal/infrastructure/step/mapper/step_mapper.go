package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/step/dto"
	"github.com/jairoprogramador/fastdeploy/internal/domain/command/values"
)

func ToDomain(commandsDTO dto.StepCommandsDTO) ([]values.CommandValue, error) {
	commandsList := make([]values.CommandValue, len(commandsDTO))

	for _, commandDto := range commandsDTO {

		commandValue, err := values.NewCommand(
			commandDto.Name,
			commandDto.Cmd,
			commandDto.Workdir,
			toOutputDomain(commandDto.Outputs),
			commandDto.Templates.Path)

		if err != nil {
			return []values.CommandValue{}, err
		}
		commandsList = append(commandsList, commandValue)
	}
	return commandsList, nil
}

func toOutputDomain(outputsDTO []dto.StepOutputDTO) []values.OutputValue {
	outputsList := make([]values.OutputValue, len(outputsDTO))
	for _, outputDto := range outputsDTO {
		validationValue := values.RegexValidation(outputDto.Probe)

		outputValue, err := values.NewOutput(
			outputDto.Name,
			outputDto.Description,
			validationValue)

		if err != nil {
			return []values.OutputValue{}
		}
		outputsList = append(outputsList, outputValue)
	}
	return outputsList
}
