package mapper

import (
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

func VarsToDTO(variables []vos.Variable) map[string]string {
	varsMap := map[string]string{}
	for _, variable := range variables {
		varsMap[variable.Key()] = variable.Value()
	}
	return varsMap
}

func VarsToDomain(varsMap map[string]string) []vos.Variable {
	variables := []vos.Variable{}
	for key, value := range varsMap {
		variable, err := vos.NewVariable(key, value)
		if err == nil {
			variables = append(variables, variable)
		}
	}
	return variables
}