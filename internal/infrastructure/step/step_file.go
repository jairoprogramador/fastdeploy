package step

import (
	"os"
	"github.com/jairoprogramador/fastdeploy/internal/domain/step/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/command/values"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/step/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/step/mapper"
	"gopkg.in/yaml.v3"
)

var DEFAULT_COMMANDS_VALUES = []values.CommandValue{}

type StepFile struct {}

func NewStepFile() port.StepPort {
	return &StepFile{}
}

func (sf *StepFile) LoadCommands(pathStepFile string) ([]values.CommandValue, error) {
	commandsFile, err := os.ReadFile(pathStepFile)
	if err != nil {
		if os.IsNotExist(err) {
			return DEFAULT_COMMANDS_VALUES, nil
		}
		return DEFAULT_COMMANDS_VALUES, err
	}

	var stepCommands dto.StepCommandsDTO
	if err := yaml.Unmarshal(commandsFile, &stepCommands); err != nil {
		return DEFAULT_COMMANDS_VALUES, err
	}

	return mapper.ToDomain(stepCommands)
}