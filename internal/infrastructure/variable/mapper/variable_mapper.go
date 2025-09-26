package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/variable/values"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/variable/dto"
)

func ToDomain(variablesDTO dto.VariablesDTO) ([]values.VariableValue, error) {
	variablesList := make([]values.VariableValue, len(variablesDTO))
	for _, variableDTO := range variablesDTO {
		variablesList = append(variablesList, values.NewVariable(variableDTO.Name, variableDTO.Value))
	}
	return variablesList, nil
}