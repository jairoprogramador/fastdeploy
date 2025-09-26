package variable

import (
	"os"

	"github.com/jairoprogramador/fastdeploy/internal/domain/variable/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/variable/values"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/variable/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/variable/mapper"
	"gopkg.in/yaml.v3"
)

var DEFAULT_VARIABLES_VALUES = []values.VariableValue{}

type VariableFile struct {}

func NewVariableFile() port.VariablePort {
	return &VariableFile{}
}

func (v *VariableFile) Load(pathFile string) ([]values.VariableValue, error) {
	data, err := os.ReadFile(pathFile)
	if err != nil {
		if os.IsNotExist(err) {
			return DEFAULT_VARIABLES_VALUES, nil
		}
		return DEFAULT_VARIABLES_VALUES, err
	}

	var result dto.VariablesDTO
	if err := yaml.Unmarshal(data, &result); err != nil {
		return DEFAULT_VARIABLES_VALUES, err
	}

	return mapper.ToDomain(result)
}